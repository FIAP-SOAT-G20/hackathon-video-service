package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

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
				s.mockController.EXPECT().List(gomock.Any(), gomock.Any(), dto.ListVideosInput{
					StatusExclude: []valueobject.VideoStatus{valueobject.CANCELLED, valueobject.COMPLETED},
					Page:          1,
					Limit:         10,
					Sort:          "status:d,created_at",
				}).Return([]byte(s.responses["list_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, res.Body.String(), s.responses["list_success"])
			},
		},
		{
			name: "success - with query",
			url:  "/videos?customer_id=1&status=OPEN,PENDING",
			setupMocks: func() {
				s.mockController.EXPECT().List(gomock.Any(), gomock.Any(), dto.ListVideosInput{
					CustomerID:    1,
					Status:        []valueobject.VideoStatus{valueobject.OPEN, valueobject.PENDING},
					StatusExclude: []valueobject.VideoStatus{valueobject.CANCELLED, valueobject.COMPLETED},
					Page:          1,
					Limit:         10,
					Sort:          "status:d,created_at",
				}).Return([]byte(s.responses["list_success_with_query"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, res.Body.String(), s.responses["list_success_with_query"])
			},
		},
		{
			name:       "invalid query - customer_id",
			url:        "/videos?customer_id=invalid",
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
		{
			name:       "invalid query - status",
			url:        "/videos?status=invalid",
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
		{
			name: "controller error",
			url:  "/videos",
			setupMocks: func() {
				s.mockController.EXPECT().List(gomock.Any(), gomock.Any(), dto.ListVideosInput{
					StatusExclude: []valueobject.VideoStatus{valueobject.CANCELLED, valueobject.COMPLETED},
					Page:          1,
					Limit:         10,
					Sort:          "status:d,created_at",
				}).Return(nil, domain.NewInternalError(nil))
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
				s.mockController.EXPECT().
					Create(gomock.Any(), gomock.Any(), dto.CreateVideoInput{CustomerID: 1}).
					Return([]byte(s.responses["create_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["create_success"])
			},
		},
		{
			name:       "invalid request - body is not a valid json",
			url:        "/videos",
			body:       strings.NewReader("invalid"),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
			},
		},
		{
			name:       "invalid request - customer_id is not a number",
			url:        "/videos",
			body:       strings.NewReader(s.requests["create_invalid_body"]),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
			},
		},
		{
			name: "controller error",
			url:  "/videos",
			body: strings.NewReader(s.requests["create_success"]),
			setupMocks: func() {
				s.mockController.EXPECT().
					Create(gomock.Any(), gomock.Any(), dto.CreateVideoInput{CustomerID: 1}).
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
				s.mockController.EXPECT().
					Get(gomock.Any(), gomock.Any(), dto.GetVideoInput{ID: 5}).
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
				s.mockController.EXPECT().
					Get(gomock.Any(), gomock.Any(), dto.GetVideoInput{ID: 5}).
					Return(nil, domain.NewNotFoundError(domain.ErrNotFound))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_not_found"])
			},
		},
		{
			name:       "invalid request - id is not a number",
			url:        "/videos/invalid",
			setupMocks: func() {},
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
				s.mockController.EXPECT().
					Update(gomock.Any(), gomock.Any(), dto.UpdateVideoInput{
						ID:         15,
						CustomerID: 5,
						Status:     valueobject.PENDING,
					}).
					Return([]byte(s.responses["update_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["update_success"])
			},
		},
		{
			name:       "invalid request - body is not a valid json",
			url:        "/videos/5",
			body:       strings.NewReader("invalid"),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_body"])
			},
		},
		{
			name:       "invalid request - customer_id is not a number",
			url:        "/videos/5",
			body:       strings.NewReader(s.requests["update_invalid_body"]),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_body"])
			},
		},
		{
			name:       "invalid request - id is not a number",
			url:        "/videos/invalid",
			body:       strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {},
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
				s.mockController.EXPECT().
					Update(gomock.Any(), gomock.Any(), dto.UpdateVideoInput{
						ID:         15,
						CustomerID: 5,
						Status:     valueobject.PENDING,
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
				s.mockController.EXPECT().
					Update(gomock.Any(), gomock.Any(), dto.UpdateVideoInput{
						ID:         15,
						CustomerID: 5,
						Status:     valueobject.PENDING,
					}).
					Return([]byte(s.responses["update_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["update_success"])
			},
		},
		{
			name:       "invalid request - body is not a valid json",
			url:        "/videos/5",
			body:       strings.NewReader("invalid"),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_body"])
			},
		},
		{
			name:       "invalid request - customer_id is not a number",
			url:        "/videos/5",
			body:       strings.NewReader(s.requests["update_invalid_body"]),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_body"])
			},
		},
		{
			name:       "invalid request - id is not a number",
			url:        "/videos/invalid",
			body:       strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {},
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
				s.mockController.EXPECT().
					Update(gomock.Any(), gomock.Any(), dto.UpdateVideoInput{
						ID:         15,
						CustomerID: 5,
						Status:     valueobject.PENDING,
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
				s.mockController.EXPECT().
					Delete(gomock.Any(), gomock.Any(), dto.DeleteVideoInput{ID: 9}).
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
				s.mockController.EXPECT().
					Delete(gomock.Any(), gomock.Any(), dto.DeleteVideoInput{ID: 9}).
					Return(nil, domain.NewNotFoundError(domain.ErrNotFound))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_not_found"])
			},
		},
		{
			name:       "invalid request - id is not a number",
			url:        "/videos/invalid",
			setupMocks: func() {},
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

			// Act
			s.router.ServeHTTP(w, req)

			// Assert
			tt.checkResult(t, w)
		})
	}
}
