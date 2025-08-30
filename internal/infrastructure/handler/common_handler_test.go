package handler_test

import (
	"maps"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/logger"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/middleware"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/server"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/util"
	"github.com/gin-gonic/gin"
)

func newRouter() *gin.Engine {
	router := gin.New()
	gin.SetMode(gin.TestMode)
	router.Use(middleware.ErrorHandler(logger.NewLogger(""))) // Remove log output
	server.RegisterCustomValidation()
	return router
}

func addCommonResponses(r *map[string]string) {
	commonResponses, err := util.ReadGoldenFiles("common",
		"error_invalid_parameter", "error_internal_error", "error_not_found",
	)
	if err != nil {
		panic(err)
	}
	maps.Copy(*r, commonResponses)
}
