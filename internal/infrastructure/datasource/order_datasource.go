package datasource

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type orderDataSource struct {
	db *gorm.DB
}

type orderKey string

func NewOrderDataSource(db *gorm.DB) port.OrderDataSource {
	return &orderDataSource{db}
}

func (ds *orderDataSource) FindByID(ctx context.Context, id uint64) (*entity.Order, error) {
	var order entity.Order
	result := ds.db.WithContext(ctx).Preload("OrderProducts.Product").First(&order, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding order: %w", result.Error)
	}
	return &order, nil
}

func (ds *orderDataSource) FindAll(ctx context.Context, filters map[string]any, sort string, page, limit int) ([]*entity.Order, int64, error) {
	var orders []*entity.Order
	var total int64

	query := ds.db.WithContext(ctx).Preload("OrderProducts.Product")

	// Apply filters
	for key, value := range filters {
		switch key {
		case "statuses":
			if statuses, ok := value.([]valueobject.OrderStatus); ok && len(statuses) > 0 {
				query = query.Where("status IN ?", statuses)
			}
		case "statuses_exclude":
			if statuses, ok := value.([]valueobject.OrderStatus); ok && len(statuses) > 0 {
				query = query.Where("status NOT IN ?", statuses)
			}
		case "customer_id":
			if customerID, ok := value.(uint64); ok && customerID != 0 {
				query = query.Where("customer_id = ?", customerID)
			}
		}
	}

	// Apply order
	if sort != "" {
		query = query.Order(sort)
	}

	// Count total before pagination
	if err := query.Model(&entity.Order{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting orders: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("error finding orders: %w", err)
	}

	return orders, total, nil
}

func (ds *orderDataSource) Create(ctx context.Context, order *entity.Order) error {
	if err := ds.db.WithContext(ctx).Create(order).Error; err != nil {
		return fmt.Errorf("error creating order: %w", err)
	}
	return nil
}

func (ds *orderDataSource) Update(ctx context.Context, order *entity.Order) error {
	result := ds.db.WithContext(ctx).Preload("OrderProducts").Save(order)
	if result.Error != nil {
		return fmt.Errorf("error updating order: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (ds *orderDataSource) Delete(ctx context.Context, id uint64) error {
	// Delete all order products first
	if err := ds.db.WithContext(ctx).Where("order_id = ?", id).Delete(&entity.OrderProduct{}).Error; err != nil {
		return fmt.Errorf("error deleting order products: %w", err)
	}

	result := ds.db.WithContext(ctx).Delete(&entity.Order{}, id)
	if result.Error != nil {
		return fmt.Errorf("error deleting order: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (ds *orderDataSource) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new context with the transaction
		keyPrincipalID := orderKey(uuid.NewString())
		txCtx := context.WithValue(ctx, keyPrincipalID, tx)
		return fn(txCtx)
	})
}
