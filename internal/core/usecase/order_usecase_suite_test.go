package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
	mockport "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port/mocks"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/usecase"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type OrderUsecaseSuiteTest struct {
	suite.Suite
	mockOrders              []*entity.Order
	mockOrderHistoryUseCase *mockport.MockOrderHistoryUseCase
	mockGateway             *mockport.MockOrderGateway
	useCase                 port.OrderUseCase
	ctx                     context.Context
}

func (s *OrderUsecaseSuiteTest) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	s.mockOrderHistoryUseCase = mockport.NewMockOrderHistoryUseCase(ctrl)
	s.mockGateway = mockport.NewMockOrderGateway(ctrl)
	s.useCase = usecase.NewOrderUseCase(s.mockGateway, s.mockOrderHistoryUseCase)
	s.ctx = context.Background()
	currentTime := time.Now()
	s.mockOrders = []*entity.Order{
		{
			ID:         1,
			CustomerID: uint64(1),
			Status:     valueobject.PENDING,
			CreatedAt:  currentTime,
			UpdatedAt:  currentTime,
		},
		{
			ID:         2,
			CustomerID: uint64(2),
			Status:     valueobject.RECEIVED,
			CreatedAt:  currentTime,
			UpdatedAt:  currentTime,
		},
	}
}

func TestOrderUsecaseSuiteTest(t *testing.T) {
	suite.Run(t, new(OrderUsecaseSuiteTest))
}
