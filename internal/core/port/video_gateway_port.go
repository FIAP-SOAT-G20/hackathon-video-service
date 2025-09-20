package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
)

type VideoGateway interface {
	FindByID(ctx context.Context, id uint64) (*entity.Video, error)
	FindAll(ctx context.Context, customerId uint64, status []valueobject.VideoStatus, statusExclude []valueobject.VideoStatus, hash string, page, limit int, sort string) ([]*entity.Video, int64, error)
	Create(ctx context.Context, video *entity.Video) error
	Update(ctx context.Context, video *entity.Video) error
	Delete(ctx context.Context, id uint64) error
}
