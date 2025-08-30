package usecase

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type categoryUseCase struct {
	gateway port.CategoryGateway
}

// NewCategoryUseCase creates a new categoryUseCase
func NewCategoryUseCase(gateway port.CategoryGateway) port.CategoryUseCase {
	return &categoryUseCase{gateway}
}

// List returns a list of Categories
func (uc *categoryUseCase) List(ctx context.Context, i dto.ListCategoriesInput) ([]*entity.Category, int64, error) {
	categories, total, err := uc.gateway.FindAll(ctx, i.Name, i.Page, i.Limit)
	if err != nil {
		return nil, 0, domain.NewInternalError(err)
	}

	return categories, total, nil
}

// Create creates a new Category
func (uc *categoryUseCase) Create(ctx context.Context, i dto.CreateCategoryInput) (*entity.Category, error) {
	category := i.ToEntity()

	if err := uc.gateway.Create(ctx, category); err != nil {
		return nil, domain.NewInternalError(err)
	}

	return category, nil
}

// Get returns a Category by ID
func (uc *categoryUseCase) Get(ctx context.Context, i dto.GetCategoryInput) (*entity.Category, error) {
	category, err := uc.gateway.FindByID(ctx, i.ID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}

	if category == nil {
		return nil, domain.NewNotFoundError(domain.ErrNotFound)
	}

	return category, nil
}

// Update updates a Category
func (uc *categoryUseCase) Update(ctx context.Context, i dto.UpdateCategoryInput) (*entity.Category, error) {
	category, err := uc.gateway.FindByID(ctx, i.ID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}
	if category == nil {
		return nil, domain.NewNotFoundError(domain.ErrNotFound)
	}

	category.Update(i.Name)

	if err := uc.gateway.Update(ctx, category); err != nil {
		return nil, domain.NewInternalError(err)
	}

	return category, nil
}

// Delete deletes a Category
func (uc *categoryUseCase) Delete(ctx context.Context, i dto.DeleteCategoryInput) (*entity.Category, error) {
	category, err := uc.gateway.FindByID(ctx, i.ID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}
	if category == nil {
		return nil, domain.NewNotFoundError(domain.ErrNotFound)
	}

	if err := uc.gateway.Delete(ctx, i.ID); err != nil {
		return nil, domain.NewInternalError(err)
	}

	return category, nil
}
