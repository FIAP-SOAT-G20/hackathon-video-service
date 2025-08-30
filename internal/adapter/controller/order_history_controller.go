package controller

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type OrderHistoryController struct {
	useCase port.OrderHistoryUseCase
}

func NewOrderHistoryController(useCase port.OrderHistoryUseCase) port.OrderHistoryController {
	return &OrderHistoryController{useCase}
}

func (c *OrderHistoryController) List(ctx context.Context, p port.Presenter, i dto.ListOrderHistoriesInput) ([]byte, error) {
	orderHistories, total, err := c.useCase.List(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{
		Total:  total,
		Page:   i.Page,
		Limit:  i.Limit,
		Result: orderHistories,
	})
}

func (c *OrderHistoryController) Create(ctx context.Context, p port.Presenter, i dto.CreateOrderHistoryInput) ([]byte, error) {
	orderHistory, err := c.useCase.Create(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: orderHistory})
}

func (c *OrderHistoryController) Get(ctx context.Context, p port.Presenter, i dto.GetOrderHistoryInput) ([]byte, error) {
	orderHistory, err := c.useCase.Get(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: orderHistory})
}

func (c *OrderHistoryController) Delete(ctx context.Context, p port.Presenter, i dto.DeleteOrderHistoryInput) ([]byte, error) {
	orderHistory, err := c.useCase.Delete(ctx, i)

	if err != nil {
		return nil, err
	}
	return p.Present(dto.PresenterInput{Result: orderHistory})
}
