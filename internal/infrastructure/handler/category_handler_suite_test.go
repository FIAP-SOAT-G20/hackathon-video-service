package handler_test

import (
	"context"
	"testing"

	mockport "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port/mocks"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/handler"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CategoryHandlerSuiteTest struct {
	suite.Suite
	handler        *handler.CategoryHandler
	router         *gin.Engine
	mockController *mockport.MockCategoryController
	ctx            context.Context
	requests       map[string]string // Fixture files
	responses      map[string]string // Golden files
}

func (s *CategoryHandlerSuiteTest) SetupTest() {
	// Create a new router
	s.router = newRouter()

	// Create a new handler
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	s.mockController = mockport.NewMockCategoryController(ctrl)
	s.handler = handler.NewCategoryHandler(s.mockController)
	s.ctx = context.Background()

	// Register routes
	s.router.GET("/categories", s.handler.List)
	s.router.POST("/categories", s.handler.Create)
	s.router.PUT("/categories/:id", s.handler.Update)
	s.router.GET("/categories/:id", s.handler.Get)
	s.router.DELETE("/categories/:id", s.handler.Delete)

	// Mock requests
	var err error
	s.requests, err = util.ReadFixtureFiles("category",
		"create_success", "create_invalid_body",
		"update_success", "update_invalid_body",
	)
	assert.NoError(s.T(), err)

	// Mock responses
	s.responses, err = util.ReadGoldenFiles("category",
		"list_success", "list_success_with_query",
		"create_success",
		"update_success",
		"get_success",
		"delete_success",
	)
	assert.NoError(s.T(), err)
	addCommonResponses(&s.responses)
}

func TestCategoryHandlerSuiteTest(t *testing.T) {
	suite.Run(t, new(CategoryHandlerSuiteTest))
}
