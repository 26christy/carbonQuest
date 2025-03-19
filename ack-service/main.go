package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/26christy/CarbonQuest/ack-service/handlers"
	"github.com/26christy/CarbonQuest/ack-service/service"
	"github.com/26christy/CarbonQuest/ack-service/storage"
	"github.com/26christy/CarbonQuest/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables from .env.ack-service
	err := godotenv.Load(".env.ack-service")
	if err != nil {
		log.Fatalf("Error loading .env.notification-service file: %v", err)
	}

	// Set up logging
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Get port from environment (default to 8082)
	port := os.Getenv("ACK_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	// Create a new Gin router
	router := gin.Default()

	// Default recovery middleware
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler(logger))

	// Initialize in-memory storage
	ackStorage := storage.NewMemoryStorage()

	// Use the default HTTP client
	httpClient := &http.Client{}

	// Initialize the ACK service
	ackService := service.NewACKService(ackStorage, httpClient)

	// Register API routes
	handlers.RegisterRoutes(router, ackService)

	// Start server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		logger.Infof("Ack Service running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server, logger)
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
