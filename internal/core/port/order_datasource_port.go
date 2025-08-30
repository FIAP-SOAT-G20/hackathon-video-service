package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
)

type OrderDataSource interface {
	FindByID(ctx context.Context, id uint64) (*entity.Order, error)
	FindAll(ctx context.Context, filters map[string]any, sort string, page, limit int) ([]*entity.Order, int64, error)
	Create(ctx context.Context, order *entity.Order) error
	Update(ctx context.Context, order *entity.Order) error
	Delete(ctx context.Context, id uint64) error
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
