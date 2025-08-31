package controller

import (
	"context"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
)

type VideoController struct {
	useCase port.VideoUseCase
}

func NewVideoController(useCase port.VideoUseCase) port.VideoController {
	return &VideoController{useCase}
}

func (c *VideoController) List(ctx context.Context, p port.Presenter, i dto.ListVideosInput) ([]byte, error) {
	videos, total, err := c.useCase.List(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{
		Total:  total,
		Page:   i.Page,
		Limit:  i.Limit,
		Result: videos,
	})
}

func (c *VideoController) Create(ctx context.Context, p port.Presenter, i dto.CreateVideoInput) ([]byte, error) {
	video, err := c.useCase.Create(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: video})
}

func (c *VideoController) Get(ctx context.Context, p port.Presenter, i dto.GetVideoInput) ([]byte, error) {
	video, err := c.useCase.Get(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: video})
}

func (c *VideoController) Update(ctx context.Context, p port.Presenter, i dto.UpdateVideoInput) ([]byte, error) {
	video, err := c.useCase.Update(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: video})
}

func (c *VideoController) Delete(ctx context.Context, p port.Presenter, i dto.DeleteVideoInput) ([]byte, error) {
	video, err := c.useCase.Delete(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{Result: video})
}
