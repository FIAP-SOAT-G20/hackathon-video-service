package gateway

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type orderProductGateway struct {
	dataSource port.OrderProductDataSource
}

func NewOrderProductGateway(dataSource port.OrderProductDataSource) port.OrderProductGateway {
	return &orderProductGateway{dataSource}
}

func (g *orderProductGateway) FindByID(ctx context.Context, orderId, productId uint64) (*entity.OrderProduct, error) {
	return g.dataSource.FindByID(ctx, orderId, productId)
}

func (g *orderProductGateway) FindAll(ctx context.Context, orderId uint64, productId uint64, page, limit int) ([]*entity.OrderProduct, int64, error) {
	filters := make(map[string]interface{})

	if orderId != 0 {
		filters["order_id"] = orderId
	}

	if productId != 0 {
		filters["product_id"] = productId
	}

	return g.dataSource.FindAll(ctx, filters, page, limit)
}

func (g *orderProductGateway) Create(ctx context.Context, orderProduct *entity.OrderProduct) error {
	return g.dataSource.Create(ctx, orderProduct)
}

func (g *orderProductGateway) Update(ctx context.Context, orderProduct *entity.OrderProduct) error {
	return g.dataSource.Update(ctx, orderProduct)
}

func (g *orderProductGateway) Delete(ctx context.Context, orderId, productId uint64) error {
	return g.dataSource.Delete(ctx, orderId, productId)
}
