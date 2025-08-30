package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

type OrderHistoryUseCase interface {
	List(ctx context.Context, input dto.ListOrderHistoriesInput) ([]*entity.OrderHistory, int64, error)
	Create(ctx context.Context, input dto.CreateOrderHistoryInput) (*entity.OrderHistory, error)
	Get(ctx context.Context, input dto.GetOrderHistoryInput) (*entity.OrderHistory, error)
	Delete(ctx context.Context, input dto.DeleteOrderHistoryInput) (*entity.OrderHistory, error)
}
