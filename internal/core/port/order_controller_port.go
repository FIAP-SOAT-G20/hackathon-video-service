package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

type OrderController interface {
	List(ctx context.Context, presenter Presenter, input dto.ListOrdersInput) ([]byte, error)
	Create(ctx context.Context, presenter Presenter, input dto.CreateOrderInput) ([]byte, error)
	Get(ctx context.Context, presenter Presenter, input dto.GetOrderInput) ([]byte, error)
	Update(ctx context.Context, presenter Presenter, input dto.UpdateOrderInput) ([]byte, error)
	Delete(ctx context.Context, presenter Presenter, input dto.DeleteOrderInput) ([]byte, error)
}
