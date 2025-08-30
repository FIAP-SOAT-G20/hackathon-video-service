package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
)

type ProductGateway interface {
	FindByID(ctx context.Context, id uint64) (*entity.Product, error)
	FindAll(ctx context.Context, name string, categoryID uint64, page, limit int) ([]*entity.Product, int64, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uint64) error
}
