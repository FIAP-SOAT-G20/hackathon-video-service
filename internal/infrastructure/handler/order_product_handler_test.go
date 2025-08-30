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

func (s *OrderProductHandlerSuiteTest) TestOrderProductHandler_List() {
	tests := []struct {
		name        string
		url         string
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/orders/products",
			setupMocks: func() {
				s.mockController.EXPECT().List(gomock.Any(), gomock.Any(), dto.ListOrderProductsInput{
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
			name: "success - with query - order_id",
			url:  "/orders/products?order_id=1",
			setupMocks: func() {
				s.mockController.EXPECT().List(gomock.Any(), gomock.Any(), dto.ListOrderProductsInput{
					OrderID: 1,
					Page:    1,
					Limit:   10,
				}).Return([]byte(s.responses["list_success_with_query"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, res.Body.String(), s.responses["list_success_with_query"])
			},
		},
		{
			name:       "invalid query - page",
			url:        "/orders/products?page=invalid",
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
		{
			name: "controller error",
			url:  "/orders/products",
			setupMocks: func() {
				s.mockController.EXPECT().List(gomock.Any(), gomock.Any(), dto.ListOrderProductsInput{
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

func (s *OrderProductHandlerSuiteTest) TestOrderProductHandler_Create() {
	tests := []struct {
		name        string
		url         string
		body        *strings.Reader
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/orders/products/1/1",
			body: strings.NewReader(s.requests["create_success"]),
			setupMocks: func() {
				s.mockController.EXPECT().
					Create(gomock.Any(), gomock.Any(), dto.CreateOrderProductInput{
						OrderID:   1,
						ProductID: 1,
						Quantity:  4,
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
			url:        "/orders/products/1/1",
			body:       strings.NewReader("invalid"),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
			},
		},
		{
			name:       "invalid request - body filed Quantity is not a number",
			url:        "/orders/products/1/1",
			body:       strings.NewReader(s.requests["create_invalid_body"]),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
			},
		},
		{
			name:       "invalid request - order_id is not a number",
			url:        "/orders/products/invalid/1",
			body:       strings.NewReader(s.requests["create_success"]),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
		{
			name: "controller error",
			url:  "/orders/products/1/1",
			body: strings.NewReader(s.requests["create_success"]),
			setupMocks: func() {
				s.mockController.EXPECT().
					Create(gomock.Any(), gomock.Any(), dto.CreateOrderProductInput{
						OrderID:   1,
						ProductID: 1,
						Quantity:  4,
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

func (s *OrderProductHandlerSuiteTest) TestOrderProductHandler_Get() {
	tests := []struct {
		name        string
		url         string
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/orders/products/1/1",
			setupMocks: func() {
				s.mockController.EXPECT().
					Get(gomock.Any(), gomock.Any(), dto.GetOrderProductInput{
						OrderID:   1,
						ProductID: 1,
					}).
					Return([]byte(s.responses["get_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["get_success"])
			},
		},
		{
			name: "not found",
			url:  "/orders/products/1/5",
			setupMocks: func() {
				s.mockController.EXPECT().
					Get(gomock.Any(), gomock.Any(), dto.GetOrderProductInput{
						OrderID:   1,
						ProductID: 5,
					}).
					Return(nil, domain.NewNotFoundError(domain.ErrNotFound))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_not_found"])
			},
		},
		{
			name:       "invalid request - order_id is not a number",
			url:        "/orders/products/invalid/1",
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

func (s *OrderProductHandlerSuiteTest) TestOrderProductHandler_Update() {
	tests := []struct {
		name        string
		url         string
		body        *strings.Reader
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success - update OrderProduct Quantity",
			url:  "/orders/products/1/1",
			body: strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {
				s.mockController.EXPECT().
					Update(gomock.Any(), gomock.Any(), dto.UpdateOrderProductInput{
						OrderID:   1,
						ProductID: 1,
						Quantity:  2,
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
			url:        "/orders/products/1/1",
			body:       strings.NewReader("invalid"),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_body"])
			},
		},
		{
			name:       "invalid request - body field Quantity is not a number",
			url:        "/orders/products/1/1",
			body:       strings.NewReader(s.requests["update_invalid_body"]),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_body"])
			},
		},
		{
			name:       "invalid request - order_id is not a number",
			url:        "/orders/products/invalid/1",
			body:       strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_invalid_parameter"])
			},
		},
		{
			name: "controller error",
			url:  "/orders/products/1/1",
			body: strings.NewReader(s.requests["update_success"]),
			setupMocks: func() {
				s.mockController.EXPECT().
					Update(gomock.Any(), gomock.Any(), dto.UpdateOrderProductInput{
						OrderID:   1,
						ProductID: 1,
						Quantity:  2,
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

func (s *OrderProductHandlerSuiteTest) TestOrderProductHandler_Delete() {
	tests := []struct {
		name        string
		url         string
		setupMocks  func()
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			url:  "/orders/products/1/1",
			setupMocks: func() {
				s.mockController.EXPECT().
					Delete(gomock.Any(), gomock.Any(), dto.DeleteOrderProductInput{
						OrderID:   1,
						ProductID: 1,
					}).
					Return([]byte(s.responses["delete_success"]), nil)
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["delete_success"])
			},
		},
		{
			name: "not found",
			url:  "/orders/products/1/5",
			setupMocks: func() {
				s.mockController.EXPECT().
					Delete(gomock.Any(), gomock.Any(), dto.DeleteOrderProductInput{
						OrderID:   1,
						ProductID: 5,
					}).
					Return(nil, domain.NewNotFoundError(domain.ErrNotFound))
			},
			checkResult: func(t *testing.T, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, res.Code)
				assert.Contains(t, util.RemoveAllSpaces(res.Body.String()), s.responses["error_not_found"])
			},
		},
		{
			name:       "invalid request - order_id is not a number",
			url:        "/orders/products/invalid/1",
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
