package handlers

import (
	"go-backend/internal/database"
	"go-backend/internal/middleware"
	"go-backend/internal/services"
	"go-backend/internal/utils"
	"go-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Router handles all routing configuration
type Router struct {
	engine     *gin.Engine
	db         *database.Database
	logger     *logger.Logger
	jwtService *utils.JWTService

	// Handlers
	userHandler   *UserHandler
	healthHandler *HealthHandler

	// Services
	userService *services.UserService
}

// NewRouter creates a new router with all dependencies
func NewRouter(db *database.Database, logger *logger.Logger, jwtService *utils.JWTService, corsOrigins []string) *Router {
	// Initialize Gin in release mode for production
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Initialize services
	userService := services.NewUserService(db.GetDB(), jwtService)

	// Initialize handlers
	userHandler := NewUserHandler(userService, logger)
	healthHandler := NewHealthHandler()

	router := &Router{
		engine:        engine,
		db:            db,
		logger:        logger,
		jwtService:    jwtService,
		userHandler:   userHandler,
		healthHandler: healthHandler,
		userService:   userService,
	}

	// Setup middleware
	router.setupMiddleware(corsOrigins)

	// Setup routes
	router.setupRoutes()

	return router
}

// setupMiddleware configures all middleware
func (r *Router) setupMiddleware(corsOrigins []string) {
	// Recovery middleware
	r.engine.Use(gin.Recovery())

	// Custom middleware
	r.engine.Use(middleware.ErrorHandlerMiddleware(r.logger))
	r.engine.Use(middleware.LoggerMiddleware(r.logger))
	r.engine.Use(middleware.CORSMiddleware(corsOrigins))
	r.engine.Use(middleware.SecurityHeadersMiddleware())
	r.engine.Use(middleware.RateLimitMiddleware())
}

// setupRoutes configures all API routes
func (r *Router) setupRoutes() {
	// Health endpoints (no auth required)
	r.engine.GET("/health", r.healthHandler.HealthCheck)
	r.engine.GET("/ready", r.healthHandler.ReadinessCheck)

	// API v1 routes
	v1 := r.engine.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.userHandler.Register)
			auth.POST("/login", r.userHandler.Login)
		}

		// Protected routes (require authentication)
		protected := v1.Group("", middleware.AuthMiddleware(r.jwtService))
		{
			// User profile routes (authenticated users)
			user := protected.Group("/user")
			{
				user.GET("/profile", r.userHandler.GetProfile)
				user.PUT("/profile", r.userHandler.UpdateUser) // Will need to extract ID from token
				user.POST("/change-password", r.userHandler.ChangePassword)
			}

			// Admin routes
			admin := protected.Group("/admin", middleware.RequireAdmin())
			{
				// User management (admin only)
				users := admin.Group("/users")
				{
					users.GET("", r.userHandler.GetUsers)
					users.GET("/:id", r.userHandler.GetUser)
					users.PUT("/:id", r.userHandler.UpdateUser)
					users.DELETE("/:id", r.userHandler.DeleteUser)
				}
			}

			// Moderator routes (admin and moderator)
			mod := protected.Group("/mod", middleware.RequireModerator())
			{
				// Add moderator-specific routes here
				mod.GET("/users", r.userHandler.GetUsers) // Moderators can view users
			}

			// Owner or admin routes (for user-specific resources)
			users := protected.Group("/users")
			{
				users.PUT("/:id", middleware.RequireOwnerOrAdmin(r.userHandler.GetUserIDFromParam), r.userHandler.UpdateUser)
			}
		}
	}

	// 404 handler
	r.engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Route not found",
		})
	})
}

// GetEngine returns the Gin engine
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// Start starts the HTTP server
func (r *Router) Start(address string) error {
	r.logger.WithField("address", address).Info("Starting HTTP server")
	return r.engine.Run(address)
}
