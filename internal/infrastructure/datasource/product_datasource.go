package datasource

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type productDataSource struct {
	db *gorm.DB
}

type productKey string

func NewProductDataSource(db *gorm.DB) port.ProductDataSource {
	return &productDataSource{db}
}

func (ds *productDataSource) FindByID(ctx context.Context, id uint64) (*entity.Product, error) {
	var product entity.Product
	result := ds.db.WithContext(ctx).First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding product: %w", result.Error)
	}
	return &product, nil
}

func (ds *productDataSource) FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entity.Product, int64, error) {
	var products []*entity.Product
	var total int64

	query := ds.db.WithContext(ctx)

	// Apply filters
	for key, value := range filters {
		switch key {
		case "name":
			if name, ok := value.(string); ok && name != "" {
				query = query.Where("name LIKE ?", "%"+name+"%")
			}
		case "category_id":
			if categoryID, ok := value.(uint64); ok && categoryID != 0 {
				query = query.Where("category_id = ?", categoryID)
			}
		}
	}

	// Count total before pagination
	if err := query.Model(&entity.Product{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting products: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("error finding products: %w", err)
	}

	return products, total, nil
}

func (ds *productDataSource) Create(ctx context.Context, product *entity.Product) error {
	if err := ds.db.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("error creating product: %w", err)
	}
	return nil
}

func (ds *productDataSource) Update(ctx context.Context, product *entity.Product) error {
	result := ds.db.WithContext(ctx).Save(product)
	if result.Error != nil {
		return fmt.Errorf("error updating product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (ds *productDataSource) Delete(ctx context.Context, id uint64) error {
	result := ds.db.WithContext(ctx).Delete(&entity.Product{}, id)
	if result.Error != nil {
		return fmt.Errorf("error deleting product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (ds *productDataSource) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new context with the transaction
		keyPrincipalID := productKey(uuid.NewString())
		txCtx := context.WithValue(ctx, keyPrincipalID, tx)
		return fn(txCtx)
	})
}
