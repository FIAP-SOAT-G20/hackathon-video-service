package gateway

import (
	"context"
	"strings"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type orderGateway struct {
	dataSource port.OrderDataSource
}

func NewOrderGateway(dataSource port.OrderDataSource) port.OrderGateway {
	return &orderGateway{dataSource}
}

func (g *orderGateway) FindByID(ctx context.Context, id uint64) (*entity.Order, error) {
	return g.dataSource.FindByID(ctx, id)
}

func (g *orderGateway) FindAll(
	ctx context.Context,
	customerId uint64,
	status []valueobject.OrderStatus,
	statusExclude []valueobject.OrderStatus,
	page,
	limit int,
	sort string,
) ([]*entity.Order, int64, error) {

	// Create filters
	filters := make(map[string]interface{})
	if customerId != 0 {
		filters["customer_id"] = customerId
	}
	if status != nil {
		filters["statuses"] = status
	}
	if statusExclude != nil {
		filters["statuses_exclude"] = statusExclude
	}

	// Create Sort "status:d,created_at" -> "status desc, created_at asc"
	sortFormatted := strings.ReplaceAll(sort, ":d", " desc")

	return g.dataSource.FindAll(ctx, filters, sortFormatted, page, limit)
}

func (g *orderGateway) Create(ctx context.Context, order *entity.Order) error {
	return g.dataSource.Create(ctx, order)
}

func (g *orderGateway) Update(ctx context.Context, order *entity.Order) error {
	return g.dataSource.Update(ctx, order)
}

func (g *orderGateway) Delete(ctx context.Context, id uint64) error {
	return g.dataSource.Delete(ctx, id)
}
