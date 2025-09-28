package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/FIAP-SOAT-G20/hackathon-video-service/docs"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/adapter/controller"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/adapter/gateway"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/usecase"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/database"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/datasource"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/handler"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/logger"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/middleware"
	awsclient "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/pkg/aws"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/pkg/aws/s3"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/route"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/server"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/service"
)

// @title						Hackathon Video Service API
// @version					1
// @description				### FIAP Hackathon - 10SOAT - G21
// @servers					[ { "url": "http://localhost:8080" }, { "url": "http://localhost:30001" } ]
// @BasePath					/api/v1
// @tag.name					videos
// @tag.description			List, create, update and delete videos
// @tag.name					health-check
// @tag.description			Health check
//
// @externalDocs.description	GitHub Repository
// @externalDocs.url			https://github.com/FIAP-SOAT-G20/hackathon-video-service
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and the access token.
func main() {
	cfg := config.LoadConfig()

	loggerInstance := logger.NewLogger(cfg.Environment)

	dbConfig, err := database.NewDatabase(cfg, loggerInstance)
	if err != nil {
		loggerInstance.Error("failed to connect to database", "error", err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := dbConfig.Close(); err != nil {
			loggerInstance.Error("failed to close database connection", "error", err.Error())
		}
	}()

	handlers := setupHandlers(dbConfig, cfg, loggerInstance)

	srv := server.NewServer(cfg, loggerInstance, handlers)
	if err := srv.Start(); err != nil {
		loggerInstance.Error("server failed to start", "error", err.Error())
		os.Exit(1)
	}
}

func setupHandlers(dbConfig *database.DatabaseConfig, cfg *config.Config, loggerInstance *logger.Logger) *route.Handlers {
	ctx := context.Background()

	// Create AWS client factory
	awsClientFactory, err := awsclient.NewClientFactory(ctx, cfg.AWSS3Region)
	if err != nil {
		panic(fmt.Sprintf("failed to create AWS client factory: %v", err))
	}

	// Create video datasource using factory
	videoDS, err := datasource.NewVideoDataSourceFactory(dbConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create video datasource: %v", err))
	}

	// Services
	jwtService := service.NewJWTService(cfg)
	// s3Service, err := service.NewS3Datasource(ctx, cfg.AWSS3BucketName, awsClientFactory)
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to create S3 service: %v", err))
	// }

	s3Client, err := s3.NewS3Client(awsClientFactory)
	if err != nil {
		panic(fmt.Sprintf("failed to create S3 client: %v", err))
	}

	// Create Cache middleware
	cacheMiddleware := middleware.NewCache(cfg, loggerInstance)

	// Gateways
	videoGateway := gateway.NewVideoGateway(videoDS)

	// Use cases
	videoUC := usecase.NewVideoUseCase(videoGateway, s3Client, cfg, loggerInstance)

	// Controllers
	videoController := controller.NewVideoController(videoUC)

	// Handlers
	videoHandler := handler.NewVideoHandler(videoController, jwtService, cacheMiddleware)
	healthCheckHandler := handler.NewHealthCheckHandler()
	redocHandler := handler.NewRedocHandler()

	handlers := &route.Handlers{
		Video:       videoHandler,
		HealthCheck: healthCheckHandler,
		Redoc:       redocHandler,
		JWTService:  jwtService,
	}

	return handlers
}
