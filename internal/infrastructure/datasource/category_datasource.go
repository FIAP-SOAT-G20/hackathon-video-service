package datasource

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type categoryDataSource struct {
	db *gorm.DB
}

func NewCategoryDataSource(db *gorm.DB) port.CategoryDataSource {
	return &categoryDataSource{db}
}

func (ds *categoryDataSource) FindByID(ctx context.Context, id uint64) (*entity.Category, error) {
	var category entity.Category
	result := ds.db.WithContext(ctx).First(&category, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding category: %w", result.Error)
	}
	return &category, nil
}

func (ds *categoryDataSource) FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entity.Category, int64, error) {
	var categorys []*entity.Category
	var total int64

	query := ds.db.WithContext(ctx)

	// Apply filters
	for key, value := range filters {
		switch key {
		case "name":
			if name, ok := value.(string); ok && name != "" {
				query = query.Where("name LIKE ?", "%"+name+"%")
			}
		}
	}

	// Count total before pagination
	if err := query.Model(&entity.Category{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting categorys: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&categorys).Error; err != nil {
		return nil, 0, fmt.Errorf("error finding categorys: %w", err)
	}

	return categorys, total, nil
}

func (ds *categoryDataSource) Create(ctx context.Context, category *entity.Category) error {
	if err := ds.db.WithContext(ctx).Create(category).Error; err != nil {
		return fmt.Errorf("error creating category: %w", err)
	}
	return nil
}

func (ds *categoryDataSource) Update(ctx context.Context, category *entity.Category) error {
	result := ds.db.WithContext(ctx).Save(category)
	if result.Error != nil {
		return fmt.Errorf("error updating category: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (ds *categoryDataSource) Delete(ctx context.Context, id uint64) error {
	result := ds.db.WithContext(ctx).Delete(&entity.Category{}, id)
	if result.Error != nil {
		return fmt.Errorf("error deleting category: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}
