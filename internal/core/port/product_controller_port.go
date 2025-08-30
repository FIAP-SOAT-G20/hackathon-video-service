package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

type ProductController interface {
	List(ctx context.Context, presenter Presenter, input dto.ListProductsInput) ([]byte, error)
	Create(ctx context.Context, presenter Presenter, input dto.CreateProductInput) ([]byte, error)
	Get(ctx context.Context, presenter Presenter, input dto.GetProductInput) ([]byte, error)
	Update(ctx context.Context, presenter Presenter, input dto.UpdateProductInput) ([]byte, error)
	Delete(ctx context.Context, presenter Presenter, input dto.DeleteProductInput) ([]byte, error)
}
