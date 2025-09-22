package route

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/handler"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/logger"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/middleware"
)

type Router struct {
	engine *gin.Engine
	logger *logger.Logger
}

func NewRouter(logger *logger.Logger, cfg *config.Config) *Router {
	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	// Global middlewares
	engine.Use(
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.ErrorHandler(logger),
		middleware.Recovery(logger),
		middleware.CORS(),
	)

	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Router{
		engine: engine,
		logger: logger,
	}
}

// RegisterRoutes configure all routes of the application
func (r *Router) RegisterRoutes(handlers *Handlers) {
	handlers.Redoc.Register(r.engine.Group("/redoc"))

	// API v1
	v1 := r.engine.Group("/api/v1")
	{
		// Videos routes with JWT authentication
		videosGroup := v1.Group("/videos")
		videosGroup.Use(middleware.JWTAuthMiddleware(handlers.JWTService))
		handlers.Video.Register(videosGroup)

		handlers.HealthCheck.Register(v1.Group("/health"))
	}
}

// Engine returns the gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// Handlers contains all handlers of the application
type Handlers struct {
	Video       *handler.VideoHandler
	JWTService  port.JWTService
	HealthCheck *handler.HealthCheckHandler
	Redoc       *handler.RedocHandler
}
