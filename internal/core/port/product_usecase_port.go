package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

type ProductUseCase interface {
	List(ctx context.Context, input dto.ListProductsInput) ([]*entity.Product, int64, error)
	Create(ctx context.Context, input dto.CreateProductInput) (*entity.Product, error)
	Get(ctx context.Context, input dto.GetProductInput) (*entity.Product, error)
	Update(ctx context.Context, input dto.UpdateProductInput) (*entity.Product, error)
	Delete(ctx context.Context, input dto.DeleteProductInput) (*entity.Product, error)
}
