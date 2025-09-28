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
	s.mockController = mockport.NewMockVideoController(ctrl)
	s.mockJWTService = mockport.NewMockJWTService(ctrl)

	// Create a simple cache middleware function for testing
	testLogger := logger.NewLogger("test")
	cacheStore := middleware.NewCacheStore(nil, testLogger)
	cacheMiddlewareFunc := cacheStore.CachePageMiddleware
	s.handler = handler.NewVideoHandler(s.mockController, s.mockJWTService, cacheMiddlewareFunc)
	s.ctx = context.Background()

	// JWT service mocks will be set up in individual test cases

	// Register routes with JWT middleware
	videoGroup := s.router.Group("/videos")
	videoGroup.Use(middleware.JWTAuth(s.mockJWTService))
	s.handler.Register(videoGroup)

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

func (s *VideoHandlerSuiteTest) BeforeTest(_, _ string) {
	// Reset mocks before each test case
	ctrl := gomock.NewController(s.T())
	s.mockController = mockport.NewMockVideoController(ctrl)
	s.mockJWTService = mockport.NewMockJWTService(ctrl)
	testLogger := logger.NewLogger("test")
	cacheStore := middleware.NewCacheStore(nil, testLogger)
	cacheMiddlewareFunc := cacheStore.CachePageMiddleware
	s.handler = handler.NewVideoHandler(s.mockController, s.mockJWTService, cacheMiddlewareFunc)

	// Create a new router for each test case
	s.router = newRouter()
	videoGroup := s.router.Group("/videos")
	videoGroup.Use(middleware.JWTAuth(s.mockJWTService))
	s.handler.Register(videoGroup)
}

func TestVideoHandlerSuiteTest(t *testing.T) {
	suite.Run(t, new(VideoHandlerSuiteTest))
}
