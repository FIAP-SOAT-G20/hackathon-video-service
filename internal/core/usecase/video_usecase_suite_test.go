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
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type VideoUsecaseSuiteTest struct {
	suite.Suite
	mockVideos        []*entity.Video
	mockGateway       *mockport.MockVideoGateway
	mockObjectStorage *mockport.MockObjectStorageDatasource
	useCase           port.VideoUseCase
	ctx               context.Context
	ctrl              *gomock.Controller
}

func (s *VideoUsecaseSuiteTest) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockGateway = mockport.NewMockVideoGateway(s.ctrl)
	s.mockObjectStorage = mockport.NewMockObjectStorageDatasource(s.ctrl)
	s.useCase = usecase.NewVideoUseCase(s.mockGateway, s.mockObjectStorage, config.LoadConfig())
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

func (s *VideoUsecaseSuiteTest) TearDownTest() {
	s.ctrl.Finish()
}

func TestVideoUsecaseSuiteTest(t *testing.T) {
	suite.Run(t, new(VideoUsecaseSuiteTest))
}
