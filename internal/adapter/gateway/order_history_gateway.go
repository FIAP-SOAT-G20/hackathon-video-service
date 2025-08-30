package gateway

import (
	"context"
	"time"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type orderHistoryGateway struct {
	dataSource port.OrderHistoryDataSource
}

func NewOrderHistoryGateway(dataSource port.OrderHistoryDataSource) port.OrderHistoryGateway {
	return &orderHistoryGateway{dataSource}
}

func (g *orderHistoryGateway) FindByID(ctx context.Context, id uint64) (*entity.OrderHistory, error) {
	return g.dataSource.FindByID(ctx, id)
}

func (g *orderHistoryGateway) FindAll(ctx context.Context, orderID uint64, status valueobject.OrderStatus, page, limit int) ([]*entity.OrderHistory, int64, error) {
	filters := make(map[string]interface{})

	if status != "" {
		filters["status"] = status.String()
	}
	if orderID != 0 {
		filters["orderID"] = orderID
	}

	return g.dataSource.FindAll(ctx, filters, page, limit)
}

func (g *orderHistoryGateway) Create(ctx context.Context, orderHistory *entity.OrderHistory) error {
	orderHistory.CreatedAt = time.Now()
	if orderHistory.StaffID != nil && *orderHistory.StaffID <= 0 {
		orderHistory.StaffID = nil
	}
	return g.dataSource.Create(ctx, orderHistory)
}

func (g *orderHistoryGateway) Delete(ctx context.Context, id uint64) error {
	return g.dataSource.Delete(ctx, id)
}
