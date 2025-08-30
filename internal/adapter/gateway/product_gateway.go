package gateway

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type productGateway struct {
	dataSource port.ProductDataSource
}

func NewProductGateway(dataSource port.ProductDataSource) port.ProductGateway {
	return &productGateway{dataSource}
}

func (g *productGateway) FindByID(ctx context.Context, id uint64) (*entity.Product, error) {
	return g.dataSource.FindByID(ctx, id)
}

func (g *productGateway) FindAll(ctx context.Context, name string, categoryID uint64, page, limit int) ([]*entity.Product, int64, error) {
	filters := make(map[string]interface{})

	if name != "" {
		filters["name"] = name
	}
	if categoryID != 0 {
		filters["category_id"] = categoryID
	}

	return g.dataSource.FindAll(ctx, filters, page, limit)
}

func (g *productGateway) Create(ctx context.Context, product *entity.Product) error {
	return g.dataSource.Create(ctx, product)
}

func (g *productGateway) Update(ctx context.Context, product *entity.Product) error {
	return g.dataSource.Update(ctx, product)
}

func (g *productGateway) Delete(ctx context.Context, id uint64) error {
	return g.dataSource.Delete(ctx, id)
}
