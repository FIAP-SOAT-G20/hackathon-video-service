package handler_test

import (
	"context"
	"testing"

	mockport "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port/mocks"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/handler"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/logger"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/middleware"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type VideoHandlerSuiteTest struct {
	suite.Suite
	handler        *handler.VideoHandler
	router         *gin.Engine
	mockController *mockport.MockVideoController
	mockJWTService *mockport.MockJWTService
	ctx            context.Context
	requests       map[string]string // Fixture files
	responses      map[string]string // Golden files
}

func (s *VideoHandlerSuiteTest) SetupTest() {
	// Create a new router
	s.router = newRouter()

	// Create a new handler
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	s.mockController = mockport.NewMockVideoController(ctrl)
	s.mockJWTService = mockport.NewMockJWTService(ctrl)

	// Create a simple cache middleware function for testing
	testLogger := logger.NewLogger("test")
	cacheStore := middleware.NewCacheStore(nil, testLogger)
	cacheMiddlewareFunc := cacheStore.CachePageMiddleware
	s.handler = handler.NewVideoHandler(s.mockController, s.mockJWTService, cacheMiddlewareFunc)
	s.ctx = context.Background()

	// Register routes
	s.router.GET("/videos", s.handler.List)
	s.router.POST("/videos", s.handler.Create)
	s.router.PUT("/videos/:id", s.handler.Update)
	s.router.PATCH("/videos/:id", s.handler.UpdatePartial)
	s.router.GET("/videos/:id", s.handler.Get)
	s.router.DELETE("/videos/:id", s.handler.Delete)
	s.router.GET("/videos/:id/processed", s.handler.Download)

	// Mock requests
	var err error
	s.requests, err = util.ReadFixtureFiles("video",
		"create_success", "create_invalid_body",
		"update_success", "update_invalid_body",
	)
	assert.NoError(s.T(), err)

	// Mock responses
	s.responses, err = util.ReadGoldenFiles("video",
		"list_success", "list_success_with_query",
		"create_success",
		"update_success",
		"get_success",
		"delete_success",
		"download_success",
	)
	assert.NoError(s.T(), err)

	addCommonResponses(&s.responses)
}

// func (s *VideoHandlerSuiteTest) BeforeTest(_, _ string) {}

func TestVideoHandlerSuiteTest(t *testing.T) {
	suite.Run(t, new(VideoHandlerSuiteTest))
}
