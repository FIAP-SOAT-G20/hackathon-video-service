package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (s *CategoryHandlerSuiteTest) TestCategoryHandler_List() {
	tests := []struct {
		name        string
		url         string
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/categories",
			setupMocks: func() {
				s.mockController.EXPECT().List(gomock.Any(), gomock.Any(), dto.ListCategoriesInput{
					Page:  1,
					Limit: 10,
				}).Return([]byte(s.responses["list_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, res.Body.String(), s.responses["list_success"])
			},
		},
		{
			name: "success - with query - name",
			url:  "/categories?name=Foods",
			setupMocks: func() {
				s.mockController.EXPECT().List(gomock.Any(), gomock.Any(), dto.ListCategoriesInput{
					Name:  "Foods",
					Page:  1,
					Limit: 10,
				}).Return([]byte(s.responses["list_success_with_query"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, res.Body.String(), s.responses["list_success_with_query"])
			},
		},
		{
			name:       "invalid query - page",
			url:        "/categories?page=invalid",
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_query"])
			},
		},
		{
			name: "controller error",
			url:  "/categories",
			setupMocks: func() {
				s.mockController.EXPECT().List(gomock.Any(), gomock.Any(), dto.ListCategoriesInput{
					Page:  1,
					Limit: 10,
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

func (s *CategoryHandlerSuiteTest) TestCategoryHandler_Get() {
	tests := []struct {
		name        string
		url         string
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/categories/1",
			setupMocks: func() {
				s.mockController.EXPECT().
					Get(gomock.Any(), gomock.Any(), dto.GetCategoryInput{ID: 1}).
					Return([]byte(s.responses["get_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["get_success"])
			},
		},
		{
			name: "not found",
			url:  "/categories/5",
			setupMocks: func() {
				s.mockController.EXPECT().
					Get(gomock.Any(), gomock.Any(), dto.GetCategoryInput{ID: 5}).
					Return(nil, domain.NewNotFoundError(domain.ErrNotFound))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_not_found"])
			},
		},
		{
			name:       "invalid request - id is not a number",
			url:        "/categories/invalid",
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
		{
			name: "controller error",
			url:  "/categories/1",
			setupMocks: func() {
				s.mockController.EXPECT().
					Get(gomock.Any(), gomock.Any(), dto.GetCategoryInput{ID: 1}).
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

			// Act
			s.router.ServeHTTP(w, req)

			// Assert
			tt.checkResult(t, w)
		})
	}
}

func (s *CategoryHandlerSuiteTest) TestCategoryHandler_Delete() {
	tests := []struct {
		name        string
		url         string
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/categories/6",
			setupMocks: func() {
				s.mockController.EXPECT().
					Delete(gomock.Any(), gomock.Any(), dto.DeleteCategoryInput{ID: 6}).
					Return([]byte(s.responses["delete_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["delete_success"])
			},
		},
		{
			name: "not found",
			url:  "/categories/6",
			setupMocks: func() {
				s.mockController.EXPECT().
					Delete(gomock.Any(), gomock.Any(), dto.DeleteCategoryInput{ID: 6}).
					Return(nil, domain.NewNotFoundError(domain.ErrNotFound))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_not_found"])
			},
		},
		{
			name:       "invalid request - id is not a number",
			url:        "/categories/invalid",
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

func (s *CategoryHandlerSuiteTest) TestCategoryHandler_Update() {
	tests := []struct {
		name        string
		url         string
		body        *strings.Reader
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success - update Category",
			url:  "/categories/5",
			body: strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {
				s.mockController.EXPECT().
					Update(gomock.Any(), gomock.Any(), dto.UpdateCategoryInput{
						ID:   5,
						Name: "Foods UPDATED",
					}).
					Return([]byte(s.responses["update_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["update_success"])
			},
		},
		{
			name:       "invalid request - id is not a number",
			url:        "/categories/invalid",
			body:       strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
		{
			name:       "invalid request - body is not a valid json",
			url:        "/categories/5",
			body:       strings.NewReader("invalid"),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_body"])
			},
		},
		{
			name:       "invalid request - body field Name is a number",
			url:        "/categories/5",
			body:       strings.NewReader(s.requests["update_invalid_body"]),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_body"])
			},
		},
		{
			name: "controller error",
			url:  "/categories/5",
			body: strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {
				s.mockController.EXPECT().
					Update(gomock.Any(), gomock.Any(), dto.UpdateCategoryInput{
						ID:   5,
						Name: "Foods UPDATED",
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

func (s *CategoryHandlerSuiteTest) TestCategoryHandler_Create() {
	tests := []struct {
		name        string
		url         string
		body        *strings.Reader
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/categories",
			body: strings.NewReader(s.requests["create_success"]),
			setupMocks: func() {
				s.mockController.EXPECT().
					Create(gomock.Any(), gomock.Any(), dto.CreateCategoryInput{
						Name: "Foods",
					}).
					Return([]byte(s.responses["create_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["create_success"])
			},
		},
		{
			name:       "invalid request - body is not a valid json",
			url:        "/categories",
			body:       strings.NewReader("invalid"),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
			},
		},
		{
			name:       "invalid request - body filed Name is a number",
			url:        "/categories",
			body:       strings.NewReader(s.requests["create_invalid_body"]),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
			},
		},
		{
			name: "controller error",
			url:  "/categories",
			body: strings.NewReader(s.requests["create_success"]),
			setupMocks: func() {
				s.mockController.EXPECT().
					Create(gomock.Any(), gomock.Any(), dto.CreateCategoryInput{
						Name: "Foods",
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

			// Act
			s.router.ServeHTTP(w, req)

			// Assert
			tt.checkResult(t, w)
		})
	}
}
