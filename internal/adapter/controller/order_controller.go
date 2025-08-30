package controller

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type OrderController struct {
	useCase port.OrderUseCase
}

func NewOrderController(useCase port.OrderUseCase) port.OrderController {
	return &OrderController{useCase}
}

func (c *OrderController) List(ctx context.Context, p port.Presenter, i dto.ListOrdersInput) ([]byte, error) {
	orders, total, err := c.useCase.List(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{
		Total:  total,
		Page:   i.Page,
		Limit:  i.Limit,
		Result: orders,
	})
}

func (c *OrderController) Create(ctx context.Context, p port.Presenter, i dto.CreateOrderInput) ([]byte, error) {
	order, err := c.useCase.Create(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: order})
}

func (c *OrderController) Get(ctx context.Context, p port.Presenter, i dto.GetOrderInput) ([]byte, error) {
	order, err := c.useCase.Get(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: order})
}

func (c *OrderController) Update(ctx context.Context, p port.Presenter, i dto.UpdateOrderInput) ([]byte, error) {
	order, err := c.useCase.Update(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: order})
}

func (c *OrderController) Delete(ctx context.Context, p port.Presenter, i dto.DeleteOrderInput) ([]byte, error) {
	order, err := c.useCase.Delete(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: order})
}
