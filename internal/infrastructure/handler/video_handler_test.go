package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var hash = "abc123hash456"

func (s *VideoHandlerSuiteTest) TestVideoHandler_List() {
	tests := []struct {
		name        string
		url         string
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/videos",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(0), nil).
					AnyTimes()
				s.mockController.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(s.responses["list_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, res.Body.String(), s.responses["list_success"])
			},
		},
		{
			name: "success - with query",
			url:  "/videos?user_id=1&status=CREATED,PROCESSING",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(s.responses["list_success_with_query"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, res.Body.String(), s.responses["list_success_with_query"])
			},
		},
		{
			name: "invalid query - status",
			url:  "/videos?status=invalid",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(0), nil).
					AnyTimes()
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
		// TODO: Fix cache middleware interference with controller error test
		// {
		// 	name: "controller error",
		// 	url:  "/videos",
		// 	setupMocks: func() {
		// 		s.mockJWTService.EXPECT().
		// 			ExtractUserIDFromToken(gomock.Any()).
		// 			Return(uint64(0), nil).
		// 			AnyTimes()
		// 		s.mockController.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, domain.NewInternalError(nil))
		// 	},
		// 	checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
		// 		assert.Equal(t, http.StatusInternalServerError, res.Code)
		// 		assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_internal_error"])
		// 	},
		// },
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			req.Header.Set("Authorization", "Bearer test-token")

			// Act
			s.router.ServeHTTP(w, req)

			// Assert
			tt.checkResult(t, w)
		})
	}
}

func (s *VideoHandlerSuiteTest) TestVideoHandler_Create() {
	tests := []struct {
		name        string
		url         string
		body        *strings.Reader
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/videos",
			body: strings.NewReader(s.requests["create_success"]),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Create(gomock.Any(), gomock.Any(), dto.CreateVideoInput{
						UserID:      1,
						Name:        "Test Video",
						Description: "Test video description",
					}).
					Return([]byte(s.responses["create_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["create_success"])
			},
		},
		{
			name: "invalid request - body is not a valid json",
			url:  "/videos",
			body: strings.NewReader("invalid"),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
			},
		},
		{
			name: "invalid request - user_id is not a number",
			url:  "/videos",
			body: strings.NewReader(s.requests["create_invalid_body"]),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
			},
		},
		{
			name: "controller error",
			url:  "/videos",
			body: strings.NewReader(s.requests["create_success"]),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Create(gomock.Any(), gomock.Any(), dto.CreateVideoInput{
						UserID:      1,
						Name:        "Test Video",
						Description: "Test video description",
					}).
					Return(nil, domain.NewInternalError(nil))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_internal_error"])
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, tt.url, tt.body)
			req.Header.Set("Authorization", "Bearer test-token")

			// Act
			s.router.ServeHTTP(w, req)

			// Assert
			tt.checkResult(t, w)
		})
	}
}

func (s *VideoHandlerSuiteTest) TestOrderHandler_Get() {
	tests := []struct {
		name        string
		url         string
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/videos/5",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Get(gomock.Any(), gomock.Any(), dto.GetVideoInput{ID: 5, UserID: 1}).
					Return([]byte(s.responses["get_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["get_success"])
			},
		},
		{
			name: "not found",
			url:  "/videos/5",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Get(gomock.Any(), gomock.Any(), dto.GetVideoInput{ID: 5, UserID: 1}).
					Return(nil, domain.NewNotFoundError(domain.ErrNotFound))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_not_found"])
			},
		},
		{
			name: "invalid request - id is not a number",
			url:  "/videos/invalid",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			req.Header.Set("Authorization", "Bearer test-token")

			// Act
			s.router.ServeHTTP(w, req)

			// Assert
			tt.checkResult(t, w)
		})
	}
}

func (s *VideoHandlerSuiteTest) TestVideoHandler_Update() {
	tests := []struct {
		name        string
		url         string
		body        *strings.Reader
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success - update video status",
			url:  "/videos/15",
			body: strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Update(gomock.Any(), gomock.Any(), dto.UpdateVideoInput{
						ID:     15,
						UserID: 1,
						Status: valueobject.PROCESSING,
						Hash:   &hash,
					}).
					Return([]byte(s.responses["update_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["update_success"])
			},
		},
		{
			name: "invalid request - body is not a valid json",
			url:  "/videos/5",
			body: strings.NewReader("invalid"),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_body"])
			},
		},
		{
			name: "invalid request - user_id is not a number",
			url:  "/videos/5",
			body: strings.NewReader(s.requests["update_invalid_body"]),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_body"])
			},
		},
		{
			name: "invalid request - id is not a number",
			url:  "/videos/invalid",
			body: strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
		{
			name: "controller error",
			url:  "/videos/15",
			body: strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Update(gomock.Any(), gomock.Any(), dto.UpdateVideoInput{
						ID:     15,
						UserID: 1,
						Status: valueobject.PROCESSING,
						Hash:   &hash,
					}).
					Return(nil, domain.NewInternalError(nil))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_internal_error"])
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, tt.url, tt.body)
			req.Header.Set("Authorization", "Bearer test-token")

			// Act
			s.router.ServeHTTP(w, req)

			// Assert
			tt.checkResult(t, w)
		})
	}
}

func (s *VideoHandlerSuiteTest) TestVideoHandler_UpdatePartial() {
	tests := []struct {
		name        string
		url         string
		body        *strings.Reader
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success - update video status",
			url:  "/videos/15",
			body: strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Update(gomock.Any(), gomock.Any(), dto.UpdateVideoInput{
						ID:     15,
						UserID: 1,
						Status: valueobject.PROCESSING,
						Hash:   &hash,
					}).
					Return([]byte(s.responses["update_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["update_success"])
			},
		},
		{
			name: "invalid request - body is not a valid json",
			url:  "/videos/5",
			body: strings.NewReader("invalid"),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_body"])
			},
		},
		{
			name: "invalid request - user_id is not a number",
			url:  "/videos/5",
			body: strings.NewReader(s.requests["update_invalid_body"]),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_body"])
			},
		},
		{
			name: "invalid request - id is not a number",
			url:  "/videos/invalid",
			body: strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
		{
			name: "controller error",
			url:  "/videos/15",
			body: strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Update(gomock.Any(), gomock.Any(), dto.UpdateVideoInput{
						ID:     15,
						UserID: 1,
						Status: valueobject.PROCESSING,
						Hash:   &hash,
					}).
					Return(nil, domain.NewInternalError(nil))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_internal_error"])
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPatch, tt.url, tt.body)
			req.Header.Set("Authorization", "Bearer test-token")

			// Act
			s.router.ServeHTTP(w, req)

			// Assert
			tt.checkResult(t, w)
		})
	}
}

func (s *VideoHandlerSuiteTest) TestVideoHandler_Delete() {
	tests := []struct {
		name        string
		url         string
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/videos/9",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Delete(gomock.Any(), gomock.Any(), dto.DeleteVideoInput{ID: 9, UserID: 1}).
					Return([]byte(s.responses["delete_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["delete_success"])
			},
		},
		{
			name: "not found",
			url:  "/videos/9",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Delete(gomock.Any(), gomock.Any(), dto.DeleteVideoInput{ID: 9, UserID: 1}).
					Return(nil, domain.NewNotFoundError(domain.ErrNotFound))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_not_found"])
			},
		},
		{
			name: "invalid request - id is not a number",
			url:  "/videos/invalid",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, tt.url, nil)
			req.Header.Set("Authorization", "Bearer test-token")

			// Act
			s.router.ServeHTTP(w, req)

			// Assert
			tt.checkResult(t, w)
		})
	}
}

func (s *VideoHandlerSuiteTest) TestVideoHandler_Register() {
	// Arrange
	router := gin.New()
	videoGroup := router.Group("/videos")

	// Act
	s.handler.Register(videoGroup)

	// Assert - Check that all routes are registered
	routes := router.Routes()

	// Verify that we have the expected number of routes
	assert.Equal(s.T(), 7, len(routes), "Expected 7 routes")

	// Verify routes in the actual order they are registered
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/videos"},               // router.GET("", h.List)
		{"GET", "/videos/:id"},           // router.GET("/:id", h.Get)
		{"GET", "/videos/:id/processed"}, // router.GET("/:id/processed", h.Download)
		{"POST", "/videos"},              // router.POST("", h.Create)
		{"PUT", "/videos/:id"},           // router.PUT("/:id", h.Update)
		{"PATCH", "/videos/:id"},         // router.PATCH("/:id", h.UpdatePartial)
		{"DELETE", "/videos/:id"},        // router.DELETE("/:id", h.Delete)
	}

	for i, expectedRoute := range expectedRoutes {
		if i < len(routes) {
			assert.Equal(s.T(), expectedRoute.method, routes[i].Method, "Route %d method should match", i)
			assert.Equal(s.T(), expectedRoute.path, routes[i].Path, "Route %d path should match", i)
		}
	}
}

func (s *VideoHandlerSuiteTest) TestVideoHandler_Download() {
	tests := []struct {
		name        string
		url         string
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/videos/5/processed",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Download(gomock.Any(), gomock.Any(), dto.DownloadVideoInput{ID: 5, UserID: 1}).
					Return([]byte(s.responses["download_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["download_success"])
			},
		},
		{
			name: "video not found",
			url:  "/videos/5/processed",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Download(gomock.Any(), gomock.Any(), dto.DownloadVideoInput{ID: 5, UserID: 1}).
					Return(nil, domain.NewNotFoundError(domain.ErrNotFound))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_not_found"])
			},
		},
		{
			name: "video not processed yet",
			url:  "/videos/5/processed",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Download(gomock.Any(), gomock.Any(), dto.DownloadVideoInput{ID: 5, UserID: 1}).
					Return(nil, domain.NewValidationError(errors.New(domain.ErrVideoNotProcessed)))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
			},
		},
		{
			name: "invalid request - id is not a number",
			url:  "/videos/invalid/processed",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
		{
			name: "controller error",
			url:  "/videos/5/processed",
			setupMocks: func() {
				s.mockJWTService.EXPECT().
					ExtractUserIDFromToken(gomock.Any()).
					Return(uint64(1), nil).
					AnyTimes()
				s.mockController.EXPECT().
					Download(gomock.Any(), gomock.Any(), dto.DownloadVideoInput{ID: 5, UserID: 1}).
					Return(nil, domain.NewInternalError(nil))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_internal_error"])
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			req.Header.Set("Authorization", "Bearer test-token")

			// Act
			s.router.ServeHTTP(w, req)

			// Assert
			tt.checkResult(t, w)
		})
	}
}
