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

type ProductUsecaseSuiteTest struct {
	suite.Suite
	mockProducts []*entity.Product
	mockGateway  *mockport.MockProductGateway
	useCase      port.ProductUseCase
	ctx          context.Context
}

func (s *ProductUsecaseSuiteTest) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	s.mockGateway = mockport.NewMockProductGateway(ctrl)
	s.useCase = usecase.NewProductUseCase(s.mockGateway)
	s.ctx = context.Background()
	currentTime := time.Now()
	s.mockProducts = []*entity.Product{
		{
			ID:          1,
			Name:        "Test Product 1",
			Description: "Description 1",
			Price:       99.99,
			CategoryID:  1,
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
		},
		{
			ID:          2,
			Name:        "Test Product 2",
			Description: "Description 2",
			Price:       199.99,
			CategoryID:  1,
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
		},
	}
}

func TestProductUsecaseSuiteTest(t *testing.T) {
	suite.Run(t, new(ProductUsecaseSuiteTest))
}
