package controller

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type categoryController struct {
	useCase port.CategoryUseCase
}

func NewCategoryController(useCase port.CategoryUseCase) port.CategoryController {
	return &categoryController{useCase}
}

func (c *categoryController) Create(ctx context.Context, p port.Presenter, i dto.CreateCategoryInput) ([]byte, error) {
	category, err := c.useCase.Create(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: category})
}

func (c *categoryController) Get(ctx context.Context, p port.Presenter, i dto.GetCategoryInput) ([]byte, error) {
	category, err := c.useCase.Get(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: category})
}

func (c *categoryController) List(ctx context.Context, p port.Presenter, i dto.ListCategoriesInput) ([]byte, error) {
	category, total, err := c.useCase.List(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{
		Total:  total,
		Page:   i.Page,
		Limit:  i.Limit,
		Result: category,
	})
}

func (c *categoryController) Update(ctx context.Context, p port.Presenter, i dto.UpdateCategoryInput) ([]byte, error) {
	category, err := c.useCase.Update(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: category})
}

func (c *categoryController) Delete(ctx context.Context, p port.Presenter, i dto.DeleteCategoryInput) ([]byte, error) {
	category, err := c.useCase.Delete(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: category})
}
