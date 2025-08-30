package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
)

type OrderProductDataSource interface {
	FindByID(ctx context.Context, orderId uint64, productId uint64) (*entity.OrderProduct, error)
	FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entity.OrderProduct, int64, error)
	Create(ctx context.Context, order *entity.OrderProduct) error
	Update(ctx context.Context, order *entity.OrderProduct) error
	Delete(ctx context.Context, orderId uint64, productId uint64) error
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
