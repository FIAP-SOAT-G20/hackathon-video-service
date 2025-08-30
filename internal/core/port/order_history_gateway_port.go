package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
)

type OrderHistoryGateway interface {
	FindByID(ctx context.Context, id uint64) (*entity.OrderHistory, error)
	FindAll(ctx context.Context, orderID uint64, status valueobject.OrderStatus, page, limit int) ([]*entity.OrderHistory, int64, error)
	Create(ctx context.Context, entity *entity.OrderHistory) error
	Delete(ctx context.Context, id uint64) error
}
