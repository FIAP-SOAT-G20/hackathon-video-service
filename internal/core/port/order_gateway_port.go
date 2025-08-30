package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
)

type OrderGateway interface {
	FindByID(ctx context.Context, id uint64) (*entity.Order, error)
	FindAll(ctx context.Context, customerId uint64, status []valueobject.OrderStatus, statusExclude []valueobject.OrderStatus, page, limit int, sort string) ([]*entity.Order, int64, error)
	Create(ctx context.Context, order *entity.Order) error
	Update(ctx context.Context, order *entity.Order) error
	Delete(ctx context.Context, id uint64) error
}
