package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/dto"
)

type VideoUseCase interface {
	List(ctx context.Context, input dto.ListVideosInput) ([]*entity.Video, int64, error)
	Create(ctx context.Context, input dto.CreateVideoInput) (*entity.Video, error)
	Get(ctx context.Context, input dto.GetVideoInput) (*entity.Video, error)
	Update(ctx context.Context, input dto.UpdateVideoInput) (*entity.Video, error)
	Delete(ctx context.Context, input dto.DeleteVideoInput) (*entity.Video, error)
}
