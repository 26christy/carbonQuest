package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/26christy/CarbonQuest/middleware"
	"github.com/26christy/CarbonQuest/notification-service/handlers"
	"github.com/26christy/CarbonQuest/notification-service/notifiers"
	"github.com/26christy/CarbonQuest/notification-service/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables from .env.notification-service
	err := godotenv.Load(".env.notification-service")
	if err != nil {
		log.Fatalf("Error loading .env.notification-service file: %v", err)
	}

	// Set up logging
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Get port from environment (default to 8081)
	port := os.Getenv("NOTIFICATION_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	// Use the default HTTP client
	httpClient := &http.Client{}

	notificationService := service.NewNotificationService(httpClient)
	// Start background scheduler
	notificationService.StartNotificationScheduler()

	// Register default notifier on startup
	registerDefaultNotifier(notificationService)

	// Setup router
	router := gin.Default()
	// Default recovery middleware
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler(logger))

	handlers.RegisterRoutes(router, notificationService)

	// Start server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		logger.Infof("Notification Service running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server, logger)
}

// registerDefaultNotifier registers a default notifier on startup
func registerDefaultNotifier(service service.NotificationService) {
	log.Println("Registering default notifier...")

	defaultNotifier, err := notifiers.CreateNotifier(os.Getenv("NOTIFIER_TYPE"), os.Getenv("NOTIFIER_PARAMS"))
	if err != nil {
		log.Fatalf("Failed to create default notifier: %v", err)
	}

	service.RegisterNotifier(defaultNotifier)
	log.Println("Default notifier registered successfully")
}

// gracefulShutdown handles clean shutdown on SIGINT or SIGTERM
func gracefulShutdown(server *http.Server, logger *logrus.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Wait for signal

	logger.Info("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown failed: %v", err)
	}

	logger.Info("Server shutdown completed")
}
