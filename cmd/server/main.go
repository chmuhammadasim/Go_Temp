package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-backend/internal/config"
	"go-backend/internal/database"
	"go-backend/internal/handlers"
	"go-backend/internal/utils"
	"go-backend/pkg/logger"

	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.NewLogger(cfg.Logging.Level, cfg.Logging.Format)
	log.Info("Starting application...")

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Error("Failed to close database connection")
		}
	}()

	// Run database migrations
	if err := db.Migrate(); err != nil {
		log.WithError(err).Fatal("Failed to run database migrations")
	}

	// Seed database with initial data
	if err := db.Seed(); err != nil {
		log.WithError(err).Fatal("Failed to seed database")
	}

	// Initialize JWT service
	jwtService := utils.NewJWTService(cfg)

	// Initialize router
	router := handlers.NewRouter(db, log, jwtService, cfg.CORS.Origins)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router.GetEngine(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithFields(logrus.Fields{
			"host": cfg.Server.Host,
			"port": cfg.Server.Port,
			"env":  cfg.Server.Env,
		}).Info("HTTP server starting")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start HTTP server")
		}
	}()

	// Setup graceful shutdown
	setupGracefulShutdown(server, log)
}

// setupGracefulShutdown handles graceful shutdown of the application
func setupGracefulShutdown(server *http.Server, log *logger.Logger) {
	// Create a channel to receive OS signals
	sigCh := make(chan os.Signal, 1)

	// Register the channel to receive specific signals
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-sigCh
	log.WithField("signal", sig.String()).Info("Received shutdown signal")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	log.Info("Attempting graceful shutdown...")
	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Graceful shutdown failed, forcing shutdown")
		os.Exit(1)
	}

	log.Info("Application shutdown completed")
}
