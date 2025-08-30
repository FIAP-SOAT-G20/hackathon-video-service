package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

func (s *OrderProductUsecaseSuiteTest) TestOrderProductsUseCase_List() {
	tests := []struct {
		name        string
		input       dto.ListOrderProductsInput
		setupMocks  func()
		checkResult func(*testing.T, []*entity.OrderProduct, int64, error)
	}{
		{
			name: "should list orderProducts successfully",
			input: dto.ListOrderProductsInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), uint64(0), 1, 10).
					Return(s.mockOrderProducts, int64(2), nil)
			},
			checkResult: func(t *testing.T, orderProducts []*entity.OrderProduct, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockOrderProducts, orderProducts)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should return error when repository fails",
			input: dto.ListOrderProductsInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), uint64(0), 1, 10).
					Return(nil, int64(0), assert.AnError)
			},
			checkResult: func(t *testing.T, orderProducts []*entity.OrderProduct, total int64, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderProducts)
				assert.Equal(t, int64(0), total)
			},
		},
		{
			name: "should filter by order id",
			input: dto.ListOrderProductsInput{
				OrderID: 1,
				Page:    1,
				Limit:   10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(1), uint64(0), 1, 10).
					Return(s.mockOrderProducts, int64(2), nil)
			},
			checkResult: func(t *testing.T, orderProducts []*entity.OrderProduct, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockOrderProducts, orderProducts)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should filter by product id",
			input: dto.ListOrderProductsInput{
				ProductID: 1,
				Page:      1,
				Limit:     10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), uint64(1), 1, 10).
					Return(s.mockOrderProducts, int64(2), nil)
			},
			checkResult: func(t *testing.T, orderProducts []*entity.OrderProduct, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockOrderProducts, orderProducts)
				assert.Equal(t, int64(2), total)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			orderProducts, total, err := s.useCase.List(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, orderProducts, total, err)
		})
	}
}

func (s *OrderProductUsecaseSuiteTest) TestOrderProductUseCase_Create() {
	tests := []struct {
		name        string
		input       dto.CreateOrderProductInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.OrderProduct, error)
	}{
		{
			name: "should create order-product successfully",
			input: dto.CreateOrderProductInput{
				OrderID:   1,
				ProductID: 1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(nil)
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, orderProduct)
				assert.Equal(t, uint64(1), orderProduct.OrderID)
				assert.Equal(t, uint64(1), orderProduct.ProductID)
			},
		},
		{
			name: "should return error when gateway fails",
			input: dto.CreateOrderProductInput{
				OrderID:   1,
				ProductID: 1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderProduct)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			orderProduct, err := s.useCase.Create(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, orderProduct, err)
		})
	}
}

func (s *OrderProductUsecaseSuiteTest) TestOrderProductUseCase_Get() {
	tests := []struct {
		name        string
		input       dto.GetOrderProductInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.OrderProduct, error)
	}{
		{
			name:  "should get orderProduct successfully",
			input: dto.GetOrderProductInput{OrderID: 1, ProductID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1), uint64(1)).
					Return(s.mockOrderProducts[0], nil)
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, orderProduct)
				assert.Equal(t, uint64(1), orderProduct.OrderID)
				assert.Equal(t, uint64(1), orderProduct.ProductID)
			},
		},
		{
			name:  "should return not found error when orderProduct doesn't exist",
			input: dto.GetOrderProductInput{OrderID: 1, ProductID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1), uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderProduct)
			},
		},
		{
			name:  "should return internal error when gateway fails",
			input: dto.GetOrderProductInput{OrderID: 1, ProductID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1), uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderProduct)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			orderProduct, err := s.useCase.Get(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, orderProduct, err)
		})
	}
}

func (s *OrderProductUsecaseSuiteTest) TestOrderProductUseCase_Update() {
	tests := []struct {
		name        string
		input       dto.UpdateOrderProductInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.OrderProduct, error)
	}{
		{
			name: "should update orderProduct successfully",
			input: dto.UpdateOrderProductInput{
				OrderID:   1,
				ProductID: 1,
				Quantity:  1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1), uint64(1)).
					Return(s.mockOrderProducts[0], nil)

				s.mockGateway.EXPECT().
					Update(s.ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, p *entity.OrderProduct) error {
						assert.Equal(s.T(), uint64(1), p.OrderID)
						assert.Equal(s.T(), uint64(1), p.ProductID)
						assert.Equal(s.T(), uint32(1), p.Quantity)
						return nil
					})
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, orderProduct)
				assert.Equal(t, uint64(1), orderProduct.OrderID)
				assert.Equal(t, uint64(1), orderProduct.ProductID)
				assert.Equal(t, uint32(1), orderProduct.Quantity)
			},
		},
		{
			name: "should return error when orderProduct not found",
			input: dto.UpdateOrderProductInput{
				OrderID:   1,
				ProductID: 1,
				Quantity:  1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1), uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderProduct)
			},
		},
		{
			name: "should return error when gateway find fails",
			input: dto.UpdateOrderProductInput{
				OrderID:   1,
				ProductID: 1,
				Quantity:  1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1), uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderProduct)
			},
		},
		{
			name: "should return error when gateway update fails",
			input: dto.UpdateOrderProductInput{
				OrderID:   1,
				ProductID: 1,
				Quantity:  1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1), uint64(1)).
					Return(s.mockOrderProducts[0], nil)

				s.mockGateway.EXPECT().
					Update(s.ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderProduct)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			orderProduct, err := s.useCase.Update(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, orderProduct, err)
		})
	}
}

func (s *OrderProductUsecaseSuiteTest) TestOrderProductUseCase_Delete() {
	tests := []struct {
		name        string
		input       dto.DeleteOrderProductInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.OrderProduct, error)
	}{
		{
			name:  "should delete orderProduct successfully",
			input: dto.DeleteOrderProductInput{OrderID: 1, ProductID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1), uint64(1)).
					Return(&entity.OrderProduct{OrderID: 1, ProductID: 1}, nil)

				s.mockGateway.EXPECT().
					Delete(s.ctx, uint64(1), uint64(1)).
					Return(nil)
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, orderProduct)
				assert.Equal(t, uint64(1), orderProduct.OrderID)
				assert.Equal(t, uint64(1), orderProduct.ProductID)
			},
		},
		{
			name:  "should return not found error when orderProduct doesn't exist",
			input: dto.DeleteOrderProductInput{OrderID: 1, ProductID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1), uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderProduct)
			},
		},
		{
			name:  "should return error when gateway fails on find",
			input: dto.DeleteOrderProductInput{OrderID: 1, ProductID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1), uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderProduct)
			},
		},
		{
			name:  "should return error when gateway fails on delete",
			input: dto.DeleteOrderProductInput{OrderID: 1, ProductID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1), uint64(1)).
					Return(&entity.OrderProduct{}, nil)

				s.mockGateway.EXPECT().
					Delete(s.ctx, uint64(1), uint64(1)).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, orderProduct *entity.OrderProduct, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderProduct)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			orderProduct, err := s.useCase.Delete(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, orderProduct, err)
		})
	}
}
