package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
)

type OrderHistoryDataSource interface {
	FindByID(ctx context.Context, id uint64) (*entity.OrderHistory, error)
	FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entity.OrderHistory, int64, error)
	Create(ctx context.Context, entity *entity.OrderHistory) error
	Delete(ctx context.Context, id uint64) error
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
