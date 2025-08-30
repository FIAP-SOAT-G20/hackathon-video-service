package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

type OrderProductController interface {
	List(ctx context.Context, presenter Presenter, input dto.ListOrderProductsInput) ([]byte, error)
	Create(ctx context.Context, presenter Presenter, input dto.CreateOrderProductInput) ([]byte, error)
	Get(ctx context.Context, presenter Presenter, input dto.GetOrderProductInput) ([]byte, error)
	Update(ctx context.Context, presenter Presenter, input dto.UpdateOrderProductInput) ([]byte, error)
	Delete(ctx context.Context, presenter Presenter, input dto.DeleteOrderProductInput) ([]byte, error)
}
