package controller_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/adapter/controller"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	mockport "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port/mocks"
)

// TODO: Add more test cenarios
func TestOrderController_ListOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mokOrdercUseCase := mockport.NewMockOrderUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewOrderController(mokOrdercUseCase)

	ctx := context.Background()
	input := dto.ListOrdersInput{
		CustomerID: 1,
		Status:     []valueobject.OrderStatus{"PENDING"},
		Page:       1,
		Limit:      10,
	}

	mockOrders := []*entity.Order{
		{
			ID:         1,
			CustomerID: 1,
			Status:     "PENDING",
		},
		{
			ID:         2,
			CustomerID: 1,
			Status:     "PENDING",
		},
	}

	mokOrdercUseCase.EXPECT().
		List(ctx, input).
		Return(mockOrders, int64(2), nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{
			Result: mockOrders,
			Total:  int64(2),
			Page:   1,
			Limit:  10,
		}).
		Return([]byte{}, nil)

	output, err := controller.List(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestOrderController_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mokOrdercUseCase := mockport.NewMockOrderUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewOrderController(mokOrdercUseCase)

	ctx := context.Background()
	input := dto.CreateOrderInput{
		CustomerID: 1,
	}

	mockOrder := &entity.Order{
		ID:         1,
		CustomerID: 1,
		Status:     "OPEN",
	}

	mokOrdercUseCase.EXPECT().
		Create(ctx, input).
		Return(mockOrder, nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{Result: mockOrder}).
		Return([]byte{}, nil)

	output, err := controller.Create(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestOrderController_GetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mokOrdercUseCase := mockport.NewMockOrderUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewOrderController(mokOrdercUseCase)

	ctx := context.Background()
	input := dto.GetOrderInput{
		ID: uint64(1),
	}

	mockOrder := &entity.Order{
		ID:         1,
		CustomerID: 1,
		Status:     "PENDING",
	}

	mokOrdercUseCase.EXPECT().
		Get(ctx, input).
		Return(mockOrder, nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{Result: mockOrder}).
		Return([]byte{}, nil)

	output, err := controller.Get(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestOrderController_UpdateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mokOrdercUseCase := mockport.NewMockOrderUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewOrderController(mokOrdercUseCase)

	ctx := context.Background()
	input := dto.UpdateOrderInput{
		ID:         uint64(1),
		CustomerID: 1,
		Status:     "OPEN",
	}

	mockOrder := &entity.Order{
		ID:         1,
		CustomerID: 1,
		Status:     "PENDING",
	}

	mokOrdercUseCase.EXPECT().
		Update(ctx, input).
		Return(mockOrder, nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{Result: mockOrder}).
		Return([]byte{}, nil)

	output, err := controller.Update(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestOrderController_DeleteOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mokOrdercUseCase := mockport.NewMockOrderUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewOrderController(mokOrdercUseCase)

	ctx := context.Background()
	input := dto.DeleteOrderInput{
		ID: uint64(1),
	}

	mockOrder := &entity.Order{
		ID:         1,
		CustomerID: 1,
		Status:     "PENDING",
	}

	mokOrdercUseCase.EXPECT().
		Delete(ctx, input).
		Return(mockOrder, nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{Result: mockOrder}).
		Return([]byte{}, nil)

	output, err := controller.Delete(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}
