package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

func (s *OrderUsecaseSuiteTest) TestOrdersUseCase_List() {
	tests := []struct {
		name        string
		input       dto.ListOrdersInput
		setupMocks  func()
		checkResult func(*testing.T, []*entity.Order, int64, error)
	}{
		{
			name: "should list orders successfully",
			input: dto.ListOrdersInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), nil, nil, 1, 10, "").
					Return(s.mockOrders, int64(2), nil)
			},
			checkResult: func(t *testing.T, orders []*entity.Order, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockOrders, orders)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should return error when repository fails",
			input: dto.ListOrdersInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), nil, nil, 1, 10, "").
					Return(nil, int64(0), assert.AnError)
			},
			checkResult: func(t *testing.T, orders []*entity.Order, total int64, err error) {
				assert.Error(t, err)
				assert.Nil(t, orders)
				assert.Equal(t, int64(0), total)
			},
		},
		{
			name: "should filter by status",
			input: dto.ListOrdersInput{
				Status: []valueobject.OrderStatus{"PENDING"},
				Page:   1,
				Limit:  10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), []valueobject.OrderStatus{"PENDING"}, nil, 1, 10, "").
					Return(s.mockOrders, int64(2), nil)
			},
			checkResult: func(t *testing.T, orders []*entity.Order, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockOrders, orders)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should filter by customer",
			input: dto.ListOrdersInput{
				CustomerID: 1,
				Page:       1,
				Limit:      10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(1), nil, nil, 1, 10, "").
					Return(s.mockOrders, int64(2), nil)
			},
			checkResult: func(t *testing.T, orders []*entity.Order, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockOrders, orders)
				assert.Equal(t, int64(2), total)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			orders, total, err := s.useCase.List(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, orders, total, err)
		})
	}
}

func (s *OrderUsecaseSuiteTest) TestOrderUseCase_Create() {
	tests := []struct {
		name        string
		input       dto.CreateOrderInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Order, error)
	}{
		{
			name: "should create order successfully",
			input: dto.CreateOrderInput{
				CustomerID: 1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(nil)
				s.mockOrderHistoryUseCase.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(&entity.OrderHistory{OrderID: 1, Status: valueobject.OPEN}, nil)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, uint64(1), order.CustomerID)
			},
		},
		{
			name: "should return error when gateway create fails",
			input: dto.CreateOrderInput{
				CustomerID: 1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name: "should return error when order history use case create fails",
			input: dto.CreateOrderInput{
				CustomerID: 1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(nil)

				s.mockOrderHistoryUseCase.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			order, err := s.useCase.Create(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, order, err)
		})
	}
}

func (s *OrderUsecaseSuiteTest) TestOrderUseCase_Get() {
	tests := []struct {
		name        string
		input       dto.GetOrderInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Order, error)
	}{
		{
			name:  "should get order successfully",
			input: dto.GetOrderInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockOrders[0], nil)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, uint64(1), order.ID)
			},
		},
		{
			name:  "should return not found error when order doesn't exist",
			input: dto.GetOrderInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name:  "should return internal error when gateway fails",
			input: dto.GetOrderInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			order, err := s.useCase.Get(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, order, err)
		})
	}
}

func (s *OrderUsecaseSuiteTest) TestOrderUseCase_Update() {
	tests := []struct {
		name        string
		input       dto.UpdateOrderInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Order, error)
	}{
		{
			name: "should update order successfully",
			input: dto.UpdateOrderInput{
				ID:         1,
				CustomerID: 1,
				Status:     valueobject.RECEIVED,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockOrders[0], nil)

				s.mockGateway.EXPECT().
					Update(s.ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, p *entity.Order) error {
						assert.Equal(s.T(), uint64(1), p.ID)
						return nil
					})

				s.mockOrderHistoryUseCase.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(&entity.OrderHistory{OrderID: 1, Status: valueobject.RECEIVED}, nil)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, valueobject.RECEIVED, order.Status)
			},
		},
		{
			name: "should return error when gateway find fails",
			input: dto.UpdateOrderInput{
				ID:         1,
				CustomerID: 1,
				Status:     valueobject.RECEIVED,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name: "should return error when order not found",
			input: dto.UpdateOrderInput{
				ID:         1,
				CustomerID: 1,
				Status:     valueobject.RECEIVED,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name: "should return error when customer id is different",
			input: dto.UpdateOrderInput{
				ID:         1,
				CustomerID: 2,
				Status:     valueobject.RECEIVED,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockOrders[0], nil)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.InvalidInputError{}, err)
			},
		},
		{
			name: "should return error when status is different and can't transition",
			input: dto.UpdateOrderInput{
				ID:         1,
				CustomerID: 1,
				Status:     valueobject.READY,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockOrders[0], nil)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.InvalidInputError{}, err)
			},
		},
		{
			name: "should return error when status is different and need staff id",
			input: dto.UpdateOrderInput{
				ID:         1,
				CustomerID: 1,
				Status:     valueobject.PREPARING,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockOrders[0], nil)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.InvalidInputError{}, err)
			},
		},
		{
			name: "should return error when gateway update fails",
			input: dto.UpdateOrderInput{
				ID:         1,
				CustomerID: 1,
				Status:     valueobject.RECEIVED,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockOrders[0], nil)

				s.mockGateway.EXPECT().
					Update(s.ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name: "should return error when status is different and order history use case create fails",
			input: dto.UpdateOrderInput{
				ID:         1,
				CustomerID: 1,
				Status:     valueobject.CANCELLED,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockOrders[0], nil)

				s.mockGateway.EXPECT().
					Update(s.ctx, gomock.Any()).
					Return(nil)

				s.mockOrderHistoryUseCase.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			order, err := s.useCase.Update(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, order, err)
		})
	}
}

func (s *OrderUsecaseSuiteTest) TestOrderUseCase_Delete() {
	tests := []struct {
		name        string
		input       dto.DeleteOrderInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Order, error)
	}{
		{
			name:  "should delete order successfully",
			input: dto.DeleteOrderInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(&entity.Order{ID: 1}, nil)

				s.mockGateway.EXPECT().
					Delete(s.ctx, uint64(1)).
					Return(nil)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, uint64(1), order.ID)
			},
		},
		{
			name:  "should return not found error when order doesn't exist",
			input: dto.DeleteOrderInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
			},
		},
		{
			name:  "should return error when gateway fails on find",
			input: dto.DeleteOrderInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name:  "should return error when gateway fails on delete",
			input: dto.DeleteOrderInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(&entity.Order{}, nil)

				s.mockGateway.EXPECT().
					Delete(s.ctx, uint64(1)).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, order *entity.Order, err error) {
				assert.Error(t, err)
				assert.Nil(t, order)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			order, err := s.useCase.Delete(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, order, err)
		})
	}
}
