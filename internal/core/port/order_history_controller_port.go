package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

type OrderHistoryController interface {
	List(ctx context.Context, presenter Presenter, input dto.ListOrderHistoriesInput) ([]byte, error)
	Create(ctx context.Context, presenter Presenter, input dto.CreateOrderHistoryInput) ([]byte, error)
	Get(ctx context.Context, presenter Presenter, input dto.GetOrderHistoryInput) ([]byte, error)
	Delete(ctx context.Context, presenter Presenter, input dto.DeleteOrderHistoryInput) ([]byte, error)
}
