package gateway

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type categoryGateway struct {
	dataSource port.CategoryDataSource
}

func NewCategoryGateway(dataSource port.CategoryDataSource) port.CategoryGateway {
	return &categoryGateway{dataSource}
}

func (g *categoryGateway) FindByID(ctx context.Context, id uint64) (*entity.Category, error) {
	return g.dataSource.FindByID(ctx, id)
}

func (g *categoryGateway) FindAll(ctx context.Context, name string, page, limit int) ([]*entity.Category, int64, error) {
	filters := make(map[string]interface{})

	if name != "" {
		filters["name"] = name
	}

	return g.dataSource.FindAll(ctx, filters, page, limit)
}

func (g *categoryGateway) Create(ctx context.Context, category *entity.Category) error {
	return g.dataSource.Create(ctx, category)
}

func (g *categoryGateway) Update(ctx context.Context, category *entity.Category) error {
	return g.dataSource.Update(ctx, category)
}

func (g *categoryGateway) Delete(ctx context.Context, id uint64) error {
	return g.dataSource.Delete(ctx, id)
}
