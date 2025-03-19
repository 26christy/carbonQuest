package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/26christy/CarbonQuest/alarm-service/handlers"
	"github.com/26christy/CarbonQuest/alarm-service/service"
	"github.com/26christy/CarbonQuest/alarm-service/storage"
	"github.com/26christy/CarbonQuest/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables from .env.alarm-service
	err := godotenv.Load(".env.alarm-service")
	if err != nil {
		log.Fatalf("Error loading .env.alarm-service file: %v", err)
	}
	// Set up logging
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Get port from environment (default to 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize dependencies
	store := storage.NewMemoryStorage()

	alarmService := service.NewAlarmService(store)
	alarmHandler := handlers.NewAlarmHandler(alarmService)

	// Setup Gin router with middleware
	router := gin.New()
	// Default recovery middleware
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler(logger))

	// Use injected handler
	handlers.RegisterRoutes(router, alarmHandler)

	// Start server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		logger.Infof("Alarm Service running on port %s", port)
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
