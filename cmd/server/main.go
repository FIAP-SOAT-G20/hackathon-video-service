package main

import (
	"os"

	_ "github.com/FIAP-SOAT-G20/tc4-order-service/docs"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/adapter/controller"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/adapter/gateway"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/usecase"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/database"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/datasource"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/handler"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/logger"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/route"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/server"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/service"
)

// @title						Fast Food API v3
// @version					1
// @description				### FIAP Tech Challenge Phase 3 - 10SOAT - G22
// @servers					[ { "url": "http://localhost:8080" }, { "url": "http://localhost:30001" } ]
// @BasePath					/api/v1
// @tag.name					sign-up
// @tag.description			Regiter a new customer
// @tag.name					sign-in
// @tag.description			Sign in to the system
// @tag.name					customers
// @tag.description			List, create, update and delete customers
// @tag.name					products
// @tag.description			List, create, update and delete products
// @tag.name					orders
// @tag.description			List, create, update and delete orders
// @tag.name					payments
// @tag.description			Process payments
// @tag.name					staffs
// @tag.description			List, create, update and delete staff
// @tag.name					health-check
// @tag.description			Health check
//
// @externalDocs.description	GitHub Repository
// @externalDocs.url			https://github.com/FIAP-SOAT-G20/tc4-order-service
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and the access token.
func main() {
	cfg := config.LoadConfig()

	loggerInstance := logger.NewLogger(cfg.Environment)

	db, err := database.NewPostgresConnection(cfg, loggerInstance)
	if err != nil {
		loggerInstance.Error("failed to connect to database", "error", err.Error())
		os.Exit(1)
	}

	if err := db.Migrate(); err != nil {
		loggerInstance.Error("failed to run migrations", "error", err.Error())
		os.Exit(1)
	}

	handlers := setupHandlers(db, cfg)

	srv := server.NewServer(cfg, loggerInstance, handlers)
	if err := srv.Start(); err != nil {
		loggerInstance.Error("server failed to start", "error", err.Error())
		os.Exit(1)
	}
}

func setupHandlers(db *database.Database, cfg *config.Config) *route.Handlers {
	// Datasources
	orderDS := datasource.NewOrderDataSource(db.DB)
	// Services
	jwtService := service.NewJWTService(cfg)

	// Gateways
	orderGateway := gateway.NewOrderGateway(orderDS)

	// Use cases
	orderUC := usecase.NewOrderUseCase(orderGateway)

	// Controllers
	orderController := controller.NewOrderController(orderUC)

	// Handlers
	orderHandler := handler.NewOrderHandler(orderController, jwtService)
	healthCheckHandler := handler.NewHealthCheckHandler()
	redocHandler := handler.NewRedocHandler()

	handlers := &route.Handlers{
		Order:       orderHandler,
		HealthCheck: healthCheckHandler,
		Redoc:       redocHandler,
	}

	return handlers
}
