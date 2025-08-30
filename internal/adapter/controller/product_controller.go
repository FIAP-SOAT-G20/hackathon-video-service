package controller

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type ProductController struct {
	useCase port.ProductUseCase
}

func NewProductController(useCase port.ProductUseCase) port.ProductController {
	return &ProductController{useCase}
}

func (c *ProductController) List(ctx context.Context, p port.Presenter, i dto.ListProductsInput) ([]byte, error) {
	products, total, err := c.useCase.List(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{
		Total:  total,
		Page:   i.Page,
		Limit:  i.Limit,
		Result: products,
	})
}

func (c *ProductController) Create(ctx context.Context, p port.Presenter, i dto.CreateProductInput) ([]byte, error) {
	product, err := c.useCase.Create(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: product})
}

func (c *ProductController) Get(ctx context.Context, p port.Presenter, i dto.GetProductInput) ([]byte, error) {
	product, err := c.useCase.Get(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: product})
}

func (c *ProductController) Update(ctx context.Context, p port.Presenter, i dto.UpdateProductInput) ([]byte, error) {
	product, err := c.useCase.Update(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: product})
}

func (c *ProductController) Delete(ctx context.Context, p port.Presenter, i dto.DeleteProductInput) ([]byte, error) {
	product, err := c.useCase.Delete(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: product})
}
