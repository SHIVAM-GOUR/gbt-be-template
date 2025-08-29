package routes

import (
	"gbt-be-template/internal/config"
	"gbt-be-template/internal/handlers"
	"gbt-be-template/internal/repository"
	"gbt-be-template/internal/services"
	"gbt-be-template/pkg/logger"
	"gbt-be-template/pkg/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// Router holds all dependencies for routing
type Router struct {
	cfg      *config.Config
	log      *logger.Logger
	db       *repository.Database
	repos    *repository.Repositories
	services *services.Services
}

// NewRouter creates a new router instance
func NewRouter(cfg *config.Config, log *logger.Logger, db *repository.Database, repos *repository.Repositories, services *services.Services) *Router {
	return &Router{
		cfg:      cfg,
		log:      log,
		db:       db,
		repos:    repos,
		services: services,
	}
}

// SetupRoutes configures all routes and middleware
func (rt *Router) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.Logging(rt.log))
	r.Use(middleware.Recovery(rt.log))
	r.Use(middleware.CORS(rt.cfg))
	r.Use(chiMiddleware.Timeout(rt.cfg.Server.GetTimeout()))

	// Initialize handlers
	userHandler := handlers.NewUserHandler(rt.services.User, rt.log)
	healthHandler := handlers.NewHealthHandler(rt.db, rt.log)

	// Health check routes (no auth required)
	r.Route("/health", func(r chi.Router) {
		r.Get("/", healthHandler.Health)
		r.Get("/ready", healthHandler.Ready)
		r.Get("/live", healthHandler.Live)
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes (no auth required)
		r.Post("/auth/login", userHandler.Login)
		r.Post("/auth/register", userHandler.Create)

		// Protected routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(rt.log, rt.cfg.JWT.Secret))

			// Protected auth routes
			r.Post("/auth/logout", userHandler.Logout)
			r.Get("/auth/profile", userHandler.Profile)

			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.List)
				r.Get("/{id}", userHandler.GetByID)
				r.Put("/{id}", userHandler.Update)
				r.Delete("/{id}", userHandler.Delete)
			})

			// Admin only routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAdmin(rt.log))

				// Admin user management
				r.Route("/admin/users", func(r chi.Router) {
					r.Post("/", userHandler.Create) // Admin can create users
				})
			})
		})
	})

	return r
}
