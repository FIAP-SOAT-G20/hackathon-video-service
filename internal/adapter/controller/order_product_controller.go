package controller

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type OrderProductController struct {
	useCase port.OrderProductUseCase
}

func NewOrderProductController(useCase port.OrderProductUseCase) port.OrderProductController {
	return &OrderProductController{useCase}
}

func (c *OrderProductController) List(ctx context.Context, p port.Presenter, i dto.ListOrderProductsInput) ([]byte, error) {
	orderProducts, total, err := c.useCase.List(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{
		Total:  total,
		Page:   i.Page,
		Limit:  i.Limit,
		Result: orderProducts,
	})
}

func (c *OrderProductController) Create(ctx context.Context, p port.Presenter, i dto.CreateOrderProductInput) ([]byte, error) {
	orderProduct, err := c.useCase.Create(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: orderProduct})
}

func (c *OrderProductController) Get(ctx context.Context, p port.Presenter, i dto.GetOrderProductInput) ([]byte, error) {
	orderProduct, err := c.useCase.Get(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: orderProduct})
}

func (c *OrderProductController) Update(ctx context.Context, p port.Presenter, i dto.UpdateOrderProductInput) ([]byte, error) {
	orderProduct, err := c.useCase.Update(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: orderProduct})
}

func (c *OrderProductController) Delete(ctx context.Context, p port.Presenter, i dto.DeleteOrderProductInput) ([]byte, error) {
	order, err := c.useCase.Delete(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: order})
}
