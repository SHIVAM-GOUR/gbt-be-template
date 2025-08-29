package main

import (
	"log"
	"os"

	"gbt-be-template/internal/config"
	"gbt-be-template/internal/server"
	"gbt-be-template/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	appLogger := logger.New(cfg.Logger.Level, cfg.Logger.Format)

	appLogger.WithFields(map[string]interface{}{
		"version": "1.0.0",
		"env":     cfg.Server.Env,
	}).Info("Starting application")

	// Create and start server
	srv, err := server.New(cfg, appLogger)
	if err != nil {
		appLogger.WithError(err).Fatal("Failed to create server")
	}

	// Start server (this blocks until shutdown)
	if err := srv.Start(); err != nil {
		appLogger.WithError(err).Fatal("Server failed")
		os.Exit(1)
	}

	appLogger.Info("Application stopped")
}
