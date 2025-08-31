package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
)

type VideoDataSource interface {
	FindByID(ctx context.Context, id uint64) (*entity.Video, error)
	FindAll(ctx context.Context, filters map[string]any, sort string, page, limit int) ([]*entity.Video, int64, error)
	Create(ctx context.Context, video *entity.Video) error
	Update(ctx context.Context, video *entity.Video) error
	Delete(ctx context.Context, id uint64) error
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
