package controller_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/adapter/controller"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/dto"
	mockport "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port/mocks"
)

// TODO: Add more test cenarios
func TestVideoController_ListVideos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mokVideocUseCase := mockport.NewMockVideoUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewVideoController(mokVideocUseCase)

	ctx := context.Background()
	input := dto.ListVideosInput{
		CustomerID: 1,
		Status:     []valueobject.VideoStatus{"PENDING"},
		Page:       1,
		Limit:      10,
	}

	mockVideos := []*entity.Video{
		{
			ID:         1,
			CustomerID: 1,
			Status:     "PENDING",
		},
		{
			ID:         2,
			CustomerID: 1,
			Status:     "PENDING",
		},
	}

	mokVideocUseCase.EXPECT().
		List(ctx, input).
		Return(mockVideos, int64(2), nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{
			Result: mockVideos,
			Total:  int64(2),
			Page:   1,
			Limit:  10,
		}).
		Return([]byte{}, nil)

	output, err := controller.List(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestVideoController_CreateVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mokVideocUseCase := mockport.NewMockVideoUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewVideoController(mokVideocUseCase)

	ctx := context.Background()
	input := dto.CreateVideoInput{
		CustomerID: 1,
	}

	mockVideo := &entity.Video{
		ID:         1,
		CustomerID: 1,
		Status:     "OPEN",
	}

	mokVideocUseCase.EXPECT().
		Create(ctx, input).
		Return(mockVideo, nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{Result: mockVideo}).
		Return([]byte{}, nil)

	output, err := controller.Create(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestVideoController_GetVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mokVideocUseCase := mockport.NewMockVideoUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewVideoController(mokVideocUseCase)

	ctx := context.Background()
	input := dto.GetVideoInput{
		ID: uint64(1),
	}

	mockVideo := &entity.Video{
		ID:         1,
		CustomerID: 1,
		Status:     "PENDING",
	}

	mokVideocUseCase.EXPECT().
		Get(ctx, input).
		Return(mockVideo, nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{Result: mockVideo}).
		Return([]byte{}, nil)

	output, err := controller.Get(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestVideoController_UpdateVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mokVideocUseCase := mockport.NewMockVideoUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewVideoController(mokVideocUseCase)

	ctx := context.Background()
	input := dto.UpdateVideoInput{
		ID:         uint64(1),
		CustomerID: 1,
		Status:     "OPEN",
	}

	mockVideo := &entity.Video{
		ID:         1,
		CustomerID: 1,
		Status:     "PENDING",
	}

	mokVideocUseCase.EXPECT().
		Update(ctx, input).
		Return(mockVideo, nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{Result: mockVideo}).
		Return([]byte{}, nil)

	output, err := controller.Update(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestVideoController_DeleteVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVideocUseCase := mockport.NewMockVideoUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	controller := controller.NewVideoController(mockVideocUseCase)

	ctx := context.Background()
	input := dto.DeleteVideoInput{
		ID: uint64(1),
	}

	mockVideo := &entity.Video{
		ID:         1,
		CustomerID: 1,
		Status:     "PENDING",
	}

	mockVideocUseCase.EXPECT().
		Delete(ctx, input).
		Return(mockVideo, nil)

	mockPresenter.EXPECT().
		Present(dto.PresenterInput{Result: mockVideo}).
		Return([]byte{}, nil)

	output, err := controller.Delete(ctx, mockPresenter, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}
