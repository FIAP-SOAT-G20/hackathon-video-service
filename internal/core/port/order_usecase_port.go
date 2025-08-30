package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

type OrderUseCase interface {
	List(ctx context.Context, input dto.ListOrdersInput) ([]*entity.Order, int64, error)
	Create(ctx context.Context, input dto.CreateOrderInput) (*entity.Order, error)
	Get(ctx context.Context, input dto.GetOrderInput) (*entity.Order, error)
	Update(ctx context.Context, input dto.UpdateOrderInput) (*entity.Order, error)
	Delete(ctx context.Context, input dto.DeleteOrderInput) (*entity.Order, error)
}
