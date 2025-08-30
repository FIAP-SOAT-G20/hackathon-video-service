package usecase_test

import (
	"testing"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (s *OrderHistoryUsecaseSuiteTest) TestOrderHistoriesUseCase_List() {
	tests := []struct {
		name        string
		input       dto.ListOrderHistoriesInput
		setupMocks  func()
		checkResult func(*testing.T, []*entity.OrderHistory, int64, error)
	}{
		{
			name: "should list staffs successfully",
			input: dto.ListOrderHistoriesInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				var status valueobject.OrderStatus
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), status, 1, 10).
					Return(s.mockOrderHistories, int64(2), nil)
			},
			checkResult: func(t *testing.T, orderHistories []*entity.OrderHistory, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockOrderHistories, orderHistories)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should return error when repository fails",
			input: dto.ListOrderHistoriesInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				var status valueobject.OrderStatus
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), status, 1, 10).
					Return(nil, int64(0), assert.AnError)
			},
			checkResult: func(t *testing.T, orderHistories []*entity.OrderHistory, total int64, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderHistories)
				assert.Equal(t, int64(0), total)
			},
		},
		{
			name: "should filter by orderID",
			input: dto.ListOrderHistoriesInput{
				OrderID: 1,
				Page:    1,
				Limit:   10,
			},
			setupMocks: func() {
				var status valueobject.OrderStatus
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(1), status, 1, 10).
					Return(s.mockOrderHistories, int64(2), nil)
			},
			checkResult: func(t *testing.T, orderHistories []*entity.OrderHistory, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockOrderHistories, orderHistories)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should filter by status",
			input: dto.ListOrderHistoriesInput{
				Status: "OPEN",
				Page:   1,
				Limit:  10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), valueobject.OPEN, 1, 10).
					Return(s.mockOrderHistories, int64(1), nil)

			},
			checkResult: func(t *testing.T, orderHistories []*entity.OrderHistory, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockOrderHistories, orderHistories)
				assert.Equal(t, int64(1), total)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			orderHistories, total, err := s.useCase.List(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, orderHistories, total, err)
		})
	}
}

func (s *OrderHistoryUsecaseSuiteTest) TestOrderHistoryUseCase_Create() {
	tests := []struct {
		name        string
		input       dto.CreateOrderHistoryInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.OrderHistory, error)
	}{
		{
			name: "should create staff successfully",
			input: dto.CreateOrderHistoryInput{
				OrderID: 1,
				Status:  "OPEN",
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(nil)
			},
			checkResult: func(t *testing.T, orderHistory *entity.OrderHistory, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, orderHistory)
				assert.Equal(t, uint64(1), orderHistory.OrderID)
				assert.Equal(t, valueobject.OPEN, orderHistory.Status)
			},
		},
		{
			name: "should return error when gateway fails",
			input: dto.CreateOrderHistoryInput{
				OrderID: 1,
				Status:  "OPEN",
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, orderHistory *entity.OrderHistory, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderHistory)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			orderHistory, err := s.useCase.Create(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, orderHistory, err)
		})
	}
}

func (s *OrderHistoryUsecaseSuiteTest) TestOrderHistoryUseCase_Get() {
	tests := []struct {
		name        string
		input       dto.GetOrderHistoryInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.OrderHistory, error)
	}{
		{
			name:  "should get order history successfully",
			input: dto.GetOrderHistoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockOrderHistories[0], nil)
			},
			checkResult: func(t *testing.T, orderHistory *entity.OrderHistory, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, orderHistory)
				assert.Equal(t, uint64(1), orderHistory.ID)
			},
		},
		{
			name:  "should return not found error when order history doesn't exist",
			input: dto.GetOrderHistoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, orderHistory *entity.OrderHistory, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderHistory)
			},
		},
		{
			name:  "should return internal error when gateway fails",
			input: dto.GetOrderHistoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, orderHistory *entity.OrderHistory, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderHistory)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			orderHistory, err := s.useCase.Get(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, orderHistory, err)
		})
	}
}

func (s *OrderHistoryUsecaseSuiteTest) TestOrderHistoryUseCase_Delete() {
	tests := []struct {
		name        string
		input       dto.DeleteOrderHistoryInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.OrderHistory, error)
	}{
		{
			name:  "should delete order history successfully",
			input: dto.DeleteOrderHistoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(&entity.OrderHistory{ID: 1}, nil)

				s.mockGateway.EXPECT().
					Delete(s.ctx, uint64(1)).
					Return(nil)
			},
			checkResult: func(t *testing.T, orderHistory *entity.OrderHistory, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, orderHistory)
				assert.Equal(t, uint64(1), orderHistory.ID)
			},
		},
		{
			name:  "should return not found error when order history doesn't exist",
			input: dto.DeleteOrderHistoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, orderHistory *entity.OrderHistory, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderHistory)
			},
		},
		{
			name:  "should return error when gateway fails on find",
			input: dto.DeleteOrderHistoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, orderHistory *entity.OrderHistory, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderHistory)
			},
		},
		{
			name:  "should return error when gateway fails on delete",
			input: dto.DeleteOrderHistoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(&entity.OrderHistory{}, nil)

				s.mockGateway.EXPECT().
					Delete(s.ctx, uint64(1)).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, orderHistory *entity.OrderHistory, err error) {
				assert.Error(t, err)
				assert.Nil(t, orderHistory)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			orderHistory, err := s.useCase.Delete(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, orderHistory, err)
		})
	}
}
