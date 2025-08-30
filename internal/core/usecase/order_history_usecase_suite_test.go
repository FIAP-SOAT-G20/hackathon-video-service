package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
	mockport "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port/mocks"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/usecase"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type OrderHistoryUsecaseSuiteTest struct {
	suite.Suite
	mockOrderHistories []*entity.OrderHistory
	mockGateway        *mockport.MockOrderHistoryGateway
	useCase            port.OrderHistoryUseCase
	ctx                context.Context
}

func (s *OrderHistoryUsecaseSuiteTest) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	s.mockGateway = mockport.NewMockOrderHistoryGateway(ctrl)
	s.useCase = usecase.NewOrderHistoryUseCase(s.mockGateway)
	s.ctx = context.Background()
	currentTime := time.Now()
	s.mockOrderHistories = []*entity.OrderHistory{
		{
			ID:        1,
			OrderID:   1,
			Status:    "OPEN",
			CreatedAt: currentTime,
		},
		{
			ID:        2,
			OrderID:   1,
			Status:    "PENDING",
			CreatedAt: currentTime,
		},
	}
}

func TestOrderHistoryUsecaseSuiteTest(t *testing.T) {
	suite.Run(t, new(OrderHistoryUsecaseSuiteTest))
}
