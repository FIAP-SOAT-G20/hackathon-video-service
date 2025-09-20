package usecase

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
)

type VideoUseCase struct {
	gateway         port.VideoGateway
	objectStorageDS port.ObjectStorageDatasource
	cacheService    port.CacheService
	config          *config.Config
}

// listCacheResult represents the cached result of a video list query
type listCacheResult struct {
	Videos []*entity.Video `json:"videos"`
	Total  int64           `json:"total"`
}

// NewVideoUseCase creates a new VideosUseCase
func NewVideoUseCase(gateway port.VideoGateway, osDS port.ObjectStorageDatasource, cacheService port.CacheService, cfg *config.Config) port.VideoUseCase {
	return &VideoUseCase{
		gateway:         gateway,
		objectStorageDS: osDS,
		cacheService:    cacheService,
		config:          cfg,
	}
}

// List returns a list of Videos
func (uc *VideoUseCase) List(ctx context.Context, i dto.ListVideosInput) ([]*entity.Video, int64, error) {
	// Generate cache key based on input parameters
	cacheKey := uc.generateListCacheKey(i)
	fmt.Println("Generated cache key:", cacheKey)

	// Try to get from cache first
	cachedResult, err := uc.getCachedListResult(ctx, cacheKey)
	if err == nil && cachedResult != nil {
		fmt.Println("Cache hit for key:", cacheKey)
		return cachedResult.Videos, cachedResult.Total, nil
	}

	// Cache miss - fetch from database
	videos, total, err := uc.gateway.FindAll(ctx, i.UserID, i.Status, i.StatusExclude, i.Hash, i.Page, i.Limit, i.Sort)
	if err != nil {
		return nil, 0, domain.NewInternalError(err)
	}

	// Cache the result for 5 minutes
	listResult := &listCacheResult{
		Videos: videos,
		Total:  total,
	}

	// Store in cache asynchronously to avoid blocking the response
	go func() {
		cacheDuration := 5 * time.Minute
		if cacheErr := uc.setCachedListResult(context.Background(), cacheKey, listResult, cacheDuration); cacheErr != nil {
			// Log the cache error but don't fail the request
			fmt.Printf("Failed to cache list result for key %s: %v\n", cacheKey, cacheErr)
		}
	}()

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
	keyWithFolder := fmt.Sprintf("%s/%s", uc.config.AWSS3BucketRawFolder, key)
	metadata := map[string]string{
		"video-id": fmt.Sprintf("%d", video.ID),
		"user-id":  fmt.Sprintf("%d", video.UserID),
	}
	presignedURL, err := uc.objectStorageDS.GeneratePresignedURL(
		ctx,
		uc.config.AWSS3BucketName,
		keyWithFolder,
		metadata,
		uc.config.AWSS3PresignedURLExpiration,
	)
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

	// Generate download key using the video hash + .zip
	key := fmt.Sprintf("%s.zip", video.Hash)
	keyWithFolder := fmt.Sprintf("%s/%s", uc.config.AWSS3BucketProcessedFolder, key)

	// Generate presigned download URL for processed video (valid for 1 hour)
	downloadURL, err := uc.objectStorageDS.GeneratePresignedDownloadURL(ctx, uc.config.AWSS3BucketName, keyWithFolder, 1*time.Hour)
	if err != nil {
		return entity.VideoProcessedDownload{}, domain.NewInternalError(err)
	}

	return entity.VideoProcessedDownload{URL: downloadURL}, nil
}

// generateListCacheKey creates a unique cache key for video list queries
func (uc *VideoUseCase) generateListCacheKey(input dto.ListVideosInput) string {
	// Convert status arrays to strings for consistent key generation
	statusStr := ""
	if len(input.Status) > 0 {
		statusVals := make([]string, len(input.Status))
		for i, status := range input.Status {
			statusVals[i] = string(status)
		}
		statusStr = strings.Join(statusVals, ",")
	}

	statusExcludeStr := ""
	if len(input.StatusExclude) > 0 {
		excludeVals := make([]string, len(input.StatusExclude))
		for i, status := range input.StatusExclude {
			excludeVals[i] = string(status)
		}
		statusExcludeStr = strings.Join(excludeVals, ",")
	}

	// Create a consistent string representation of the input
	keyData := fmt.Sprintf("videos:list:uid:%d:status:%s:exclude:%s:hash:%s:page:%d:limit:%d:sort:%s",
		input.UserID, statusStr, statusExcludeStr, input.Hash, input.Page, input.Limit, input.Sort)

	// Generate MD5 hash for shorter, consistent keys
	hash := md5.Sum([]byte(keyData))
	return fmt.Sprintf("videos:list:%x", hash)
}

// getCachedListResult retrieves cached video list result
func (uc *VideoUseCase) getCachedListResult(ctx context.Context, key string) (*listCacheResult, error) {
	if uc.cacheService == nil {
		fmt.Println("Cache service not available")
		return nil, fmt.Errorf("cache service not available")
	}

	// Add timeout context to prevent hanging on cache operations
	cacheCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	cachedData, err := uc.cacheService.Get(cacheCtx, key)
	if err != nil {
		fmt.Printf("Cache miss or error for key %s: %v\n", key, err)
		return nil, err
	}

	if cachedData == nil {
		fmt.Println("Cache miss for key:", key)
		return nil, fmt.Errorf("cache miss")
	}

	// Convert cached data to JSON string then unmarshal to struct
	jsonData, ok := cachedData.(string)
	if !ok {
		fmt.Println("Invalid cached data type")
		return nil, fmt.Errorf("invalid cached data type")
	}

	var result listCacheResult
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		fmt.Println("Error unmarshaling cached data:", err)
		return nil, err
	}

	return &result, nil
}

// setCachedListResult stores video list result in cache
func (uc *VideoUseCase) setCachedListResult(ctx context.Context, key string, result *listCacheResult, expiration time.Duration) error {
	if uc.cacheService == nil {
		return fmt.Errorf("cache service not available")
	}

	// Add timeout context to prevent hanging on cache operations
	cacheCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Marshal to JSON for storage
	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return uc.cacheService.Set(cacheCtx, key, string(jsonData), expiration)
}
