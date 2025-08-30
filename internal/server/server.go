package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gbt-be-template/internal/config"
	"gbt-be-template/internal/repository"
	"gbt-be-template/internal/routes"
	"gbt-be-template/internal/services"
	"gbt-be-template/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// Server represents the HTTP server
type Server struct {
	cfg    *config.Config
	log    *logger.Logger
	db     *repository.Database
	router *chi.Mux
	server *http.Server
}

// New creates a new server instance
func New(cfg *config.Config, log *logger.Logger) (*Server, error) {
	// Initialize database
	db, err := repository.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Run auto migration only in development mode when not using Docker
	// In Docker, we use proper migrations via migrate container
	skipAutoMigrate := os.Getenv("SKIP_AUTO_MIGRATE")
	log.Info("Auto migration check", "skip_auto_migrate", skipAutoMigrate, "is_development", cfg.IsDevelopment())

	if cfg.IsDevelopment() && skipAutoMigrate != "true" {
		log.Info("Running auto migration")
		if err := db.AutoMigrate(); err != nil {
			return nil, fmt.Errorf("failed to run auto migration: %w", err)
		}
	} else {
		log.Info("Skipping auto migration", "reason", "skip_auto_migrate=true or not development")
	}

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	authService := services.NewAuthService(repos.User, cfg, log)
	userService := services.NewUserService(repos.User, authService, cfg, log)

	services := &services.Services{
		User: userService,
		Auth: authService,
	}

	// Initialize router
	router := routes.NewRouter(cfg, log, db, repos, services)
	mux := router.SetupRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		cfg:    cfg,
		log:    log,
		db:     db,
		router: mux,
		server: server,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		s.log.WithFields(map[string]interface{}{
			"addr": s.server.Addr,
			"env":  s.cfg.Server.Env,
		}).Info("Starting HTTP server")

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	<-quit
	s.log.Info("Shutting down server...")

	// Graceful shutdown
	return s.Shutdown()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.WithError(err).Error("Failed to shutdown server gracefully")
		return err
	}

	// Close database connection
	if err := s.db.Close(); err != nil {
		s.log.WithError(err).Error("Failed to close database connection")
		return err
	}

	s.log.Info("Server shutdown completed")
	return nil
}

// GetAddr returns the server address
func (s *Server) GetAddr() string {
	return s.server.Addr
}

// GetRouter returns the router instance
func (s *Server) GetRouter() *chi.Mux {
	return s.router
}
