package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/dto"
)

type VideoController interface {
	List(ctx context.Context, presenter Presenter, input dto.ListVideosInput) ([]byte, error)
	Create(ctx context.Context, presenter Presenter, input dto.CreateVideoInput) ([]byte, error)
	Get(ctx context.Context, presenter Presenter, input dto.GetVideoInput) ([]byte, error)
	Update(ctx context.Context, presenter Presenter, input dto.UpdateVideoInput) ([]byte, error)
	Delete(ctx context.Context, presenter Presenter, input dto.DeleteVideoInput) ([]byte, error)
}
