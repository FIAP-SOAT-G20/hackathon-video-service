package controller_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/adapter/controller"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	mockport "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port/mocks"
)

// TODO: Add more test cenarios
func TestOrderProductController_ListOrderProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderProductUseCase := mockport.NewMockOrderProductUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewOrderProductController(mockOrderProductUseCase)

	ctx := context.Background()
	input := dto.ListOrderProductsInput{
		OrderID: 1,
		Page:    1,
		Limit:   10,
	}

	mockOrderProducts := []*entity.OrderProduct{
		{
			OrderID:   1,
			ProductID: 1,
			Quantity:  1,
		},
		{
			OrderID:   1,
			ProductID: 2,
			Quantity:  2,
		},
	}

	mockOrderProductUseCase.EXPECT().
		List(ctx, input).
		Return(mockOrderProducts, int64(2), nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{
			Result: mockOrderProducts,
			Total:  int64(2),
			Page:   1,
			Limit:  10,
		}).
		Return([]byte{}, nil)

	output, err := controller.List(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestOrderProductController_CreateOrderProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderProductUseCase := mockport.NewMockOrderProductUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewOrderProductController(mockOrderProductUseCase)

	ctx := context.Background()
	input := dto.CreateOrderProductInput{
		OrderID:   1,
		ProductID: 1,
		Quantity:  1,
	}

	mockOrderProduct := &entity.OrderProduct{
		OrderID:   1,
		ProductID: 1,
		Quantity:  1,
	}

	mockOrderProductUseCase.EXPECT().
		Create(ctx, input).
		Return(mockOrderProduct, nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{Result: mockOrderProduct}).
		Return([]byte{}, nil)

	output, err := controller.Create(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestOrderProductController_GetOrderProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderProductUseCase := mockport.NewMockOrderProductUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewOrderProductController(mockOrderProductUseCase)

	ctx := context.Background()
	input := dto.GetOrderProductInput{
		OrderID:   1,
		ProductID: 1,
	}

	mockOrderProduct := &entity.OrderProduct{
		OrderID:   1,
		ProductID: 1,
		Quantity:  1,
	}

	mockOrderProductUseCase.EXPECT().
		Get(ctx, input).
		Return(mockOrderProduct, nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{Result: mockOrderProduct}).
		Return([]byte{}, nil)

	output, err := controller.Get(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestOrderProductController_UpdateOrderProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderProductUseCase := mockport.NewMockOrderProductUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewOrderProductController(mockOrderProductUseCase)

	ctx := context.Background()
	input := dto.UpdateOrderProductInput{
		OrderID:   1,
		ProductID: 1,
		Quantity:  2,
	}

	mockOrderProduct := &entity.OrderProduct{
		OrderID:   1,
		ProductID: 1,
		Quantity:  2,
	}

	mockOrderProductUseCase.EXPECT().
		Update(ctx, input).
		Return(mockOrderProduct, nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{Result: mockOrderProduct}).
		Return([]byte{}, nil)

	output, err := controller.Update(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestOrderProductController_DeleteOrderProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderProductUseCase := mockport.NewMockOrderProductUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewOrderProductController(mockOrderProductUseCase)

	ctx := context.Background()
	input := dto.DeleteOrderProductInput{
		OrderID:   1,
		ProductID: 1,
	}

	mockOrderProduct := &entity.OrderProduct{
		OrderID:   1,
		ProductID: 1,
		Quantity:  1,
	}

	mockOrderProductUseCase.EXPECT().
		Delete(ctx, input).
		Return(mockOrderProduct, nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{Result: mockOrderProduct}).
		Return([]byte{}, nil)

	output, err := controller.Delete(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}
