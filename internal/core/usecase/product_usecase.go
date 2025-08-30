package usecase

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type productUseCase struct {
	gateway port.ProductGateway
}

// NewProductUseCase creates a new StaffUseCase
func NewProductUseCase(gateway port.ProductGateway) port.ProductUseCase {
	return &productUseCase{gateway: gateway}
}

// List returns a list of products
func (uc *productUseCase) List(ctx context.Context, i dto.ListProductsInput) ([]*entity.Product, int64, error) {
	products, total, err := uc.gateway.FindAll(ctx, i.Name, i.CategoryID, i.Page, i.Limit)
	if err != nil {
		return nil, 0, domain.NewInternalError(err)
	}

	return products, total, nil
}

// Create creates a new product
func (uc *productUseCase) Create(ctx context.Context, i dto.CreateProductInput) (*entity.Product, error) {
	product := i.ToEntity()

	if err := uc.gateway.Create(ctx, product); err != nil {
		return nil, domain.NewInternalError(err)
	}

	return product, nil
}

// Get returns a product by ID
func (uc *productUseCase) Get(ctx context.Context, i dto.GetProductInput) (*entity.Product, error) {
	staff, err := uc.gateway.FindByID(ctx, i.ID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}

	if staff == nil {
		return nil, domain.NewNotFoundError(domain.ErrNotFound)
	}

	return staff, nil
}

// Update updates a product
func (uc *productUseCase) Update(ctx context.Context, i dto.UpdateProductInput) (*entity.Product, error) {
	product, err := uc.gateway.FindByID(ctx, i.ID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}
	if product == nil {
		return nil, domain.NewNotFoundError(domain.ErrNotFound)
	}

	product.Update(i.Name, i.Description, i.Price, i.CategoryID)

	if err := uc.gateway.Update(ctx, product); err != nil {
		return nil, domain.NewInternalError(err)
	}

	return product, nil
}

// Delete deletes a product
func (uc *productUseCase) Delete(ctx context.Context, i dto.DeleteProductInput) (*entity.Product, error) {
	product, err := uc.gateway.FindByID(ctx, i.ID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}
	if product == nil {
		return nil, domain.NewNotFoundError(domain.ErrNotFound)
	}

	if err := uc.gateway.Delete(ctx, i.ID); err != nil {
		return nil, domain.NewInternalError(err)
	}

	return product, nil
}
