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

type OrderProductUsecaseSuiteTest struct {
	suite.Suite
	mockOrderProducts []*entity.OrderProduct
	mockGateway       *mockport.MockOrderProductGateway
	useCase           port.OrderProductUseCase
	ctx               context.Context
}

func (s *OrderProductUsecaseSuiteTest) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	s.mockGateway = mockport.NewMockOrderProductGateway(ctrl)
	s.useCase = usecase.NewOrderProductUseCase(s.mockGateway)
	s.ctx = context.Background()
	currentTime := time.Now()
	s.mockOrderProducts = []*entity.OrderProduct{
		{
			OrderID:   1,
			ProductID: 1,
			Quantity:  1,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
		{
			OrderID:   2,
			ProductID: 2,
			Quantity:  2,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
	}
}

func TestOrderProductUsecaseSuiteTest(t *testing.T) {
	suite.Run(t, new(OrderProductUsecaseSuiteTest))
}
