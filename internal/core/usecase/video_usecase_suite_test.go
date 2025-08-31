package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	mockport "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port/mocks"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/usecase"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type VideoUsecaseSuiteTest struct {
	suite.Suite
	mockVideos  []*entity.Video
	mockGateway *mockport.MockVideoGateway
	useCase     port.VideoUseCase
	ctx         context.Context
}

func (s *VideoUsecaseSuiteTest) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	s.mockGateway = mockport.NewMockVideoGateway(ctrl)
	s.useCase = usecase.NewVideoUseCase(s.mockGateway)
	s.ctx = context.Background()
	currentTime := time.Now()
	s.mockVideos = []*entity.Video{
		{
			ID:        1,
			UserID:    uint64(1),
			Status:    valueobject.CREATED,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
		{
			ID:        2,
			UserID:    uint64(2),
			Status:    valueobject.PROCESSING,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
	}
}

func TestVideoUsecaseSuiteTest(t *testing.T) {
	suite.Run(t, new(VideoUsecaseSuiteTest))
}
