package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
)

type VideoUseCase struct {
	gateway   port.VideoGateway
	s3Service port.S3Service
}

// NewVideoUseCase creates a new VideosUseCase
func NewVideoUseCase(gateway port.VideoGateway, s3Service port.S3Service) port.VideoUseCase {
	return &VideoUseCase{
		gateway:   gateway,
		s3Service: s3Service,
	}
}

// List returns a list of Videos
func (uc *VideoUseCase) List(ctx context.Context, i dto.ListVideosInput) ([]*entity.Video, int64, error) {
	videos, total, err := uc.gateway.FindAll(ctx, i.UserID, i.Status, i.StatusExclude, i.Hash, i.Page, i.Limit, i.Sort)
	if err != nil {
		return nil, 0, domain.NewInternalError(err)
	}

	return videos, total, nil
}

// Create creates a new Video
func (uc *VideoUseCase) Create(ctx context.Context, i dto.CreateVideoInput) (*entity.Video, error) {

	video := &entity.Video{
		UserID:      i.UserID,
		Status:      valueobject.CREATED,
		Name:        i.Name,
		Description: i.Description,
	}

	if err := uc.gateway.Create(ctx, video); err != nil {
		return nil, domain.NewInternalError(err)
	}

	// Generate S3 presigned URL for video upload
	key := fmt.Sprintf("%d.mp4", video.ID)
	presignedURL, err := uc.s3Service.GeneratePresignedURL(ctx, key, 15*time.Minute)
	if err != nil {
		// Log error but don't fail the video creation
		// The presigned URL generation is not critical for video entity creation
		presignedURL = ""
	}

	video.PresignedURL = presignedURL

	return video, nil
}

// Get returns a Video by ID
func (uc *VideoUseCase) Get(ctx context.Context, i dto.GetVideoInput) (*entity.Video, error) {
	video, err := uc.gateway.FindByID(ctx, i.ID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}

	if video == nil {
		return nil, domain.NewNotFoundError(domain.ErrNotFound)
	}

	return video, nil
}

// Update updates a Video
func (uc *VideoUseCase) Update(ctx context.Context, i dto.UpdateVideoInput) (*entity.Video, error) {
	video, err := uc.gateway.FindByID(ctx, i.ID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}

	if video == nil {
		return nil, domain.NewNotFoundError(domain.ErrNotFound)
	}

	if i.UserID != 0 && video.UserID != i.UserID {
		return nil, domain.NewInvalidInputError(domain.ErrInvalidBody)
	}

	statusHasChanged := video.Status != i.Status
	if i.Status != "" && statusHasChanged {
		if !valueobject.StatusCanTransitionTo(video.Status, i.Status) {
			return nil, domain.NewInvalidInputError(domain.ErrVideoInvalidStatusTransition)
		}
	}

	video.Update(i.UserID, i.Status, i.Name, i.Description, i.Hash)

	if err := uc.gateway.Update(ctx, video); err != nil {
		return nil, domain.NewInternalError(err)
	}

	return video, nil
}

// Delete deletes a Video
func (uc *VideoUseCase) Delete(ctx context.Context, i dto.DeleteVideoInput) (*entity.Video, error) {
	video, err := uc.gateway.FindByID(ctx, i.ID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}
	if video == nil {
		return nil, domain.NewNotFoundError(domain.ErrNotFound)
	}

	if err := uc.gateway.Delete(ctx, i.ID); err != nil {
		return nil, domain.NewInternalError(err)
	}

	return video, nil
}

// Download generates a presigned URL for downloading a processed video
func (uc *VideoUseCase) Download(ctx context.Context, i dto.DownloadVideoInput) (entity.VideoProcessedDownload, error) {
	video, err := uc.gateway.FindByID(ctx, i.ID)
	if err != nil {
		return entity.VideoProcessedDownload{}, domain.NewInternalError(err)
	}

	if video == nil {
		return entity.VideoProcessedDownload{}, domain.NewNotFoundError(domain.ErrNotFound)
	}

	// Only allow download for FINISHED videos
	if video.Status != valueobject.FINISHED {
		return entity.VideoProcessedDownload{}, domain.NewInvalidInputError(domain.ErrVideoNotProcessed)
	}

	// Generate download key using the video hash
	key := video.Hash

	// Generate presigned download URL for processed video (valid for 1 hour)
	downloadURL, err := uc.s3Service.GeneratePresignedDownloadURL(ctx, key, 1*time.Hour)
	if err != nil {
		return entity.VideoProcessedDownload{}, domain.NewInternalError(err)
	}

	return entity.VideoProcessedDownload{URL: downloadURL}, nil
}
