package datasource

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
)

type videoDataSource struct {
	db *gorm.DB
}

type videoKey string

func NewVideoDataSource(db *gorm.DB) port.VideoDataSource {
	return &videoDataSource{db}
}

func (ds *videoDataSource) FindByID(ctx context.Context, id uint64) (*entity.Video, error) {
	var video entity.Video
	result := ds.db.WithContext(ctx).First(&video, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding video: %w", result.Error)
	}
	return &video, nil
}

func (ds *videoDataSource) FindAll(ctx context.Context, filters map[string]any, sort string, page, limit int) ([]*entity.Video, int64, error) {
	var videos []*entity.Video
	var total int64

	query := ds.db.WithContext(ctx)

	// Apply filters
	for key, value := range filters {
		switch key {
		case "statuses":
			if statuses, ok := value.([]valueobject.VideoStatus); ok && len(statuses) > 0 {
				query = query.Where("status IN ?", statuses)
			}
		case "statuses_exclude":
			if statuses, ok := value.([]valueobject.VideoStatus); ok && len(statuses) > 0 {
				query = query.Where("status NOT IN ?", statuses)
			}
		case "user_id":
			if customerID, ok := value.(uint64); ok && customerID != 0 {
				query = query.Where("user_id = ?", customerID)
			}
		case "hash":
			if hash, ok := value.(string); ok && hash != "" {
				query = query.Where("hash = ?", hash)
			}
		}
	}

	// Apply video sorting
	if sort != "" {
		query = query.Order(sort)
	}

	// Count total before pagination
	if err := query.Model(&entity.Video{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting videos: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&videos).Error; err != nil {
		return nil, 0, fmt.Errorf("error finding videos: %w", err)
	}

	return videos, total, nil
}

func (ds *videoDataSource) Create(ctx context.Context, video *entity.Video) error {
	if err := ds.db.WithContext(ctx).Create(video).Error; err != nil {
		return fmt.Errorf("error creating video: %w", err)
	}
	return nil
}

func (ds *videoDataSource) Update(ctx context.Context, video *entity.Video) error {
	result := ds.db.WithContext(ctx).Save(video)
	if result.Error != nil {
		return fmt.Errorf("error updating video: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (ds *videoDataSource) Delete(ctx context.Context, id uint64) error {
	result := ds.db.WithContext(ctx).Delete(&entity.Video{}, id)
	if result.Error != nil {
		return fmt.Errorf("error deleting video: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (ds *videoDataSource) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new context with the transaction
		keyPrincipalID := videoKey(uuid.NewString())
		txCtx := context.WithValue(ctx, keyPrincipalID, tx)
		return fn(txCtx)
	})
}
