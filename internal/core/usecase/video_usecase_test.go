package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/dto"
)

func (s *VideoUsecaseSuiteTest) TestVideosUseCase_List() {
	tests := []struct {
		name        string
		input       dto.ListVideosInput
		setupMocks  func()
		checkResult func(*testing.T, []*entity.Video, int64, error)
	}{
		{
			name: "should list videos successfully",
			input: dto.ListVideosInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), nil, nil, "", 1, 10, "").
					Return(s.mockVideos, int64(2), nil)
			},
			checkResult: func(t *testing.T, videos []*entity.Video, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockVideos, videos)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should return error when repository fails",
			input: dto.ListVideosInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), nil, nil, "", 1, 10, "").
					Return(nil, int64(0), assert.AnError)
			},
			checkResult: func(t *testing.T, videos []*entity.Video, total int64, err error) {
				assert.Error(t, err)
				assert.Nil(t, videos)
				assert.Equal(t, int64(0), total)
			},
		},
		{
			name: "should filter by status",
			input: dto.ListVideosInput{
				Status: []valueobject.VideoStatus{"PENDING"},
				Page:   1,
				Limit:  10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), []valueobject.VideoStatus{"PENDING"}, nil, "", 1, 10, "").
					Return(s.mockVideos, int64(2), nil)
			},
			checkResult: func(t *testing.T, videos []*entity.Video, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockVideos, videos)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should filter by customer",
			input: dto.ListVideosInput{
				UserID: 1,
				Page:   1,
				Limit:  10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(1), nil, nil, "", 1, 10, "").
					Return(s.mockVideos, int64(2), nil)
			},
			checkResult: func(t *testing.T, videos []*entity.Video, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockVideos, videos)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should filter by hash",
			input: dto.ListVideosInput{
				Hash:  "abc123hash456",
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, uint64(0), nil, nil, "abc123hash456", 1, 10, "").
					Return(s.mockVideos, int64(1), nil)
			},
			checkResult: func(t *testing.T, videos []*entity.Video, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockVideos, videos)
				assert.Equal(t, int64(1), total)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			videos, total, err := s.useCase.List(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, videos, total, err)
		})
	}
}

func (s *VideoUsecaseSuiteTest) TestVideoUseCase_Create() {
	tests := []struct {
		name        string
		input       dto.CreateVideoInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Video, error)
	}{
		{
			name: "should create video successfully",
			input: dto.CreateVideoInput{
				UserID: 1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(nil)
				s.mockS3Service.EXPECT().
					GeneratePresignedURL(s.ctx, gomock.Any(), gomock.Any()).
					Return("https://s3.example.com/presigned-url", nil)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, video)
				assert.Equal(t, uint64(1), video.UserID)
				assert.NotEmpty(t, video.PresignedURL)
			},
		},
		{
			name: "should return error when gateway create fails",
			input: dto.CreateVideoInput{
				UserID: 1,
			},
			setupMocks: func() {
				s.mockS3Service.EXPECT().
					GeneratePresignedURL(s.ctx, gomock.Any(), gomock.Any()).
					Return("https://s3.example.com/presigned-url", nil)
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.Error(t, err)
				assert.Nil(t, video)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name: "should create video even when S3 service fails",
			input: dto.CreateVideoInput{
				UserID: 1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(nil)
				s.mockS3Service.EXPECT().
					GeneratePresignedURL(s.ctx, gomock.Any(), gomock.Any()).
					Return("", assert.AnError)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, video)
				assert.Equal(t, uint64(1), video.UserID)
				assert.Empty(t, video.PresignedURL) // Should be empty when S3 fails
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			video, err := s.useCase.Create(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, video, err)
		})
	}
}

func (s *VideoUsecaseSuiteTest) TestVideoUseCase_Get() {
	tests := []struct {
		name        string
		input       dto.GetVideoInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Video, error)
	}{
		{
			name:  "should get video successfully",
			input: dto.GetVideoInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockVideos[0], nil)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, video)
				assert.Equal(t, uint64(1), video.ID)
			},
		},
		{
			name:  "should return not found error when video doesn't exist",
			input: dto.GetVideoInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.Error(t, err)
				assert.Nil(t, video)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name:  "should return internal error when gateway fails",
			input: dto.GetVideoInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.Error(t, err)
				assert.Nil(t, video)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			video, err := s.useCase.Get(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, video, err)
		})
	}
}

func (s *VideoUsecaseSuiteTest) TestVideoUseCase_Update() {
	tests := []struct {
		name        string
		input       dto.UpdateVideoInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Video, error)
	}{
		{
			name: "should update video successfully",
			input: dto.UpdateVideoInput{
				ID:     1,
				UserID: 1,
				Status: valueobject.PROCESSING,
			},
			setupMocks: func() {
				// Create a fresh video entity for this test to avoid state pollution
				video := &entity.Video{
					ID:        1,
					UserID:    uint64(1),
					Status:    valueobject.CREATED,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(video, nil)

				s.mockGateway.EXPECT().
					Update(s.ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, p *entity.Video) error {
						assert.Equal(s.T(), uint64(1), p.ID)
						return nil
					})
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, video)
				assert.Equal(t, valueobject.PROCESSING, video.Status)
			},
		},
		{
			name: "should return error when gateway find fails",
			input: dto.UpdateVideoInput{
				ID:     1,
				UserID: 1,
				Status: valueobject.PROCESSING,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.Error(t, err)
				assert.Nil(t, video)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name: "should return error when video not found",
			input: dto.UpdateVideoInput{
				ID:     1,
				UserID: 1,
				Status: valueobject.PROCESSING,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.Error(t, err)
				assert.Nil(t, video)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name: "should return error when customer id is different",
			input: dto.UpdateVideoInput{
				ID:     1,
				UserID: 2,
				Status: valueobject.PROCESSING,
			},
			setupMocks: func() {
				// Create a fresh video entity for this test to avoid state pollution
				video := &entity.Video{
					ID:        1,
					UserID:    uint64(1), // Different from input UserID to test error case
					Status:    valueobject.CREATED,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(video, nil)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.Error(t, err)
				assert.Nil(t, video)
				assert.IsType(t, &domain.InvalidInputError{}, err)
			},
		},
		{
			name: "should return error when status is different and can't transition",
			input: dto.UpdateVideoInput{
				ID:     1,
				UserID: 1,
				Status: valueobject.FINISHED,
			},
			setupMocks: func() {
				// Create a fresh video entity for this test to avoid state pollution
				video := &entity.Video{
					ID:        1,
					UserID:    uint64(1),
					Status:    valueobject.CREATED, // Cannot transition directly from CREATED to FINISHED
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(video, nil)
				// No expectation for Update because it should not be called due to invalid transition
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.Error(t, err)
				assert.Nil(t, video)
				assert.IsType(t, &domain.InvalidInputError{}, err)
			},
		},
		{
			name: "should return error when gateway update fails",
			input: dto.UpdateVideoInput{
				ID:     1,
				UserID: 1,
				Status: valueobject.PROCESSING,
			},
			setupMocks: func() {
				// Create a fresh video entity for this test to avoid state pollution
				video := &entity.Video{
					ID:        1,
					UserID:    uint64(1),
					Status:    valueobject.CREATED,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(video, nil)

				s.mockGateway.EXPECT().
					Update(s.ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.Error(t, err)
				assert.Nil(t, video)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name: "should update video successfully when transitioning from CREATED to FAILED",
			input: dto.UpdateVideoInput{
				ID:     1,
				UserID: 1,
				Status: valueobject.FAILED,
			},
			setupMocks: func() {
				// Create a fresh video entity for this test to avoid state pollution
				video := &entity.Video{
					ID:        1,
					UserID:    uint64(1),
					Status:    valueobject.CREATED,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(video, nil)

				s.mockGateway.EXPECT().
					Update(s.ctx, gomock.Any()).
					Return(nil)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, video)
				assert.Equal(t, valueobject.FAILED, video.Status)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			video, err := s.useCase.Update(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, video, err)
		})
	}
}

func (s *VideoUsecaseSuiteTest) TestVideoUseCase_Delete() {
	tests := []struct {
		name        string
		input       dto.DeleteVideoInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Video, error)
	}{
		{
			name:  "should delete video successfully",
			input: dto.DeleteVideoInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(&entity.Video{ID: 1}, nil)

				s.mockGateway.EXPECT().
					Delete(s.ctx, uint64(1)).
					Return(nil)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, video)
				assert.Equal(t, uint64(1), video.ID)
			},
		},
		{
			name:  "should return not found error when video doesn't exist",
			input: dto.DeleteVideoInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.Error(t, err)
				assert.Nil(t, video)
			},
		},
		{
			name:  "should return error when gateway fails on find",
			input: dto.DeleteVideoInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.Error(t, err)
				assert.Nil(t, video)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name:  "should return error when gateway fails on delete",
			input: dto.DeleteVideoInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(&entity.Video{}, nil)

				s.mockGateway.EXPECT().
					Delete(s.ctx, uint64(1)).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, video *entity.Video, err error) {
				assert.Error(t, err)
				assert.Nil(t, video)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			video, err := s.useCase.Delete(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, video, err)
		})
	}
}
