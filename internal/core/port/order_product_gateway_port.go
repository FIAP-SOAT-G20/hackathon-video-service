package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
)

type OrderProductGateway interface {
	FindByID(ctx context.Context, orderId uint64, productId uint64) (*entity.OrderProduct, error)
	FindAll(ctx context.Context, orderId uint64, productId uint64, page, limit int) ([]*entity.OrderProduct, int64, error)
	Create(ctx context.Context, orderProduct *entity.OrderProduct) error
	Update(ctx context.Context, orderProduct *entity.OrderProduct) error
	Delete(ctx context.Context, orderId uint64, productId uint64) error
}
