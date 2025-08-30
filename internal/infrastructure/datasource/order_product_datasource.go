package datasource

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type orderProductDataSource struct {
	db *gorm.DB
}

type orderProductKey string

func NewOrderProductDataSource(db *gorm.DB) port.OrderProductDataSource {
	return &orderProductDataSource{db}
}

func (ds *orderProductDataSource) FindByID(ctx context.Context, orderId, productId uint64) (*entity.OrderProduct, error) {
	var orderProduct entity.OrderProduct
	result := ds.db.WithContext(ctx).Preload("Order").Preload("Product").First(&orderProduct, "order_id = ? AND product_id = ?", orderId, productId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding orderProduct: %w", result.Error)
	}
	return &orderProduct, nil
}

func (ds *orderProductDataSource) FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entity.OrderProduct, int64, error) {
	var orderProducts []*entity.OrderProduct
	var total int64

	query := ds.db.WithContext(ctx).Preload("Order").Preload("Product")

	// Apply filters
	for key, value := range filters {
		switch key {
		case "order_id":
			if orderID, ok := value.(uint64); ok && orderID != 0 {
				query = query.Where("order_id = ?", orderID)
			}
		}
	}

	// Count total before pagination
	if err := query.Model(&entity.OrderProduct{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting orderProducts: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&orderProducts).Error; err != nil {
		return nil, 0, fmt.Errorf("error finding orderProducts: %w", err)
	}

	return orderProducts, total, nil
}

func (ds *orderProductDataSource) Create(ctx context.Context, orderProduct *entity.OrderProduct) error {
	if err := ds.db.WithContext(ctx).Create(orderProduct).Error; err != nil {
		return fmt.Errorf("error creating orderProduct: %w", err)
	}

	// Preload related entities
	if err := ds.db.WithContext(ctx).Preload("Order").Preload("Product").First(orderProduct, "order_id = ? AND product_id = ?", orderProduct.OrderID, orderProduct.ProductID).Error; err != nil {
		return fmt.Errorf("error preloading orderProduct: %w", err)
	}

	return nil
}

func (ds *orderProductDataSource) Update(ctx context.Context, orderProduct *entity.OrderProduct) error {
	result := ds.db.WithContext(ctx).Model(orderProduct).Where("order_id = ? AND product_id = ?", orderProduct.OrderID, orderProduct.ProductID).Updates(orderProduct)
	if result.Error != nil {
		return fmt.Errorf("error updating orderProduct: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (ds *orderProductDataSource) Delete(ctx context.Context, orderId, productId uint64) error {
	result := ds.db.WithContext(ctx).Delete(&entity.OrderProduct{}, "order_id = ? AND product_id = ?", orderId, productId)
	if result.Error != nil {
		return fmt.Errorf("error deleting orderProduct: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (ds *orderProductDataSource) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new context with the transaction
		keyPrincipalID := orderProductKey(uuid.NewString())
		txCtx := context.WithValue(ctx, keyPrincipalID, tx)
		return fn(txCtx)
	})
}
