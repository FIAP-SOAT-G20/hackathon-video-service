package datasource

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type orderHistoryDataSource struct {
	db *gorm.DB
}

type orderHistoryKey string

func NewOrderHistoryDataSource(db *gorm.DB) port.OrderHistoryDataSource {
	return &orderHistoryDataSource{
		db: db,
	}
}

func (ds *orderHistoryDataSource) FindByID(ctx context.Context, id uint64) (*entity.OrderHistory, error) {
	var orderHistory entity.OrderHistory
	result := ds.db.WithContext(ctx).First(&orderHistory, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding orderHistory: %w", result.Error)
	}
	return &orderHistory, nil
}

func (ds *orderHistoryDataSource) FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entity.OrderHistory, int64, error) {
	var orderHistories []*entity.OrderHistory
	var total int64

	query := ds.db.WithContext(ctx)

	// Apply filters
	for key, value := range filters {
		switch key {
		case "orderID":
			if orderID, ok := value.(uint64); ok && orderID != 0 {
				query = query.Where("order_id = ?", orderID)
			}
		case "status":
			if status, ok := value.(string); ok && status != "" && status != "UNDEFINED" {
				query = query.Where("status LIKE ?", status)
			}
		}

	}

	// Count total before pagination
	if err := query.Model(&entity.OrderHistory{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting orderHistorys: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&orderHistories).Error; err != nil {
		return nil, 0, fmt.Errorf("error finding orderHistorys: %w", err)
	}

	return orderHistories, total, nil
}

func (ds *orderHistoryDataSource) Create(ctx context.Context, orderHistory *entity.OrderHistory) error {
	if err := ds.db.WithContext(ctx).Create(orderHistory).Error; err != nil {
		return fmt.Errorf("error creating orderHistory: %w", err)
	}
	return nil
}

func (ds *orderHistoryDataSource) Update(ctx context.Context, orderHistory *entity.OrderHistory) error {
	result := ds.db.WithContext(ctx).Save(orderHistory)
	if result.Error != nil {
		return fmt.Errorf("error updating orderHistory: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (ds *orderHistoryDataSource) Delete(ctx context.Context, id uint64) error {
	result := ds.db.WithContext(ctx).Delete(&entity.OrderHistory{}, id)
	if result.Error != nil {
		return fmt.Errorf("error deleting orderHistory: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (ds *orderHistoryDataSource) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new context with the transaction
		keyPrincipalID := orderHistoryKey(uuid.NewString())
		txCtx := context.WithValue(ctx, keyPrincipalID, tx)
		return fn(txCtx)
	})
}
