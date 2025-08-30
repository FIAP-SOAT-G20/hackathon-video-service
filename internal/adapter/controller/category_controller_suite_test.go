package controller_test

import (
	"context"
	"testing"
	"time"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/adapter/controller"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/adapter/presenter"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
	mockport "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CategoryControllerSuiteTest struct {
	suite.Suite
	mockUseCase    *mockport.MockCategoryUseCase
	mockPresenter  port.Presenter
	mockCategory   *entity.Category
	mockCategories []*entity.Category
	controller     port.CategoryController
	ctx            context.Context
}

func (s *CategoryControllerSuiteTest) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	s.mockUseCase = mockport.NewMockCategoryUseCase(ctrl)
	s.mockPresenter = presenter.NewCategoryJsonPresenter()
	mockDateAt2, _ := time.Parse(time.RFC3339, "2025-03-06T17:03:28-03:00")
	mockDateAt3, _ := time.Parse(time.RFC3339, "2025-03-06T17:03:58-03:00")
	s.mockCategory = &entity.Category{
		ID:        6,
		Name:      "Foods",
		CreatedAt: mockDateAt2,
		UpdatedAt: mockDateAt3,
	}
	mockDateAt, _ := time.Parse(time.RFC3339, "2025-02-28T16:28:18Z")
	s.mockCategories = []*entity.Category{
		{
			ID:        1,
			Name:      "Foods",
			CreatedAt: mockDateAt,
			UpdatedAt: mockDateAt,
		},
		{
			ID:        2,
			Name:      "Beverages",
			CreatedAt: mockDateAt,
			UpdatedAt: mockDateAt,
		},
	}
	s.controller = controller.NewCategoryController(s.mockUseCase)
	s.ctx = context.Background()
}

func TestCategoryControllerSuiteTest(t *testing.T) {
	suite.Run(t, new(CategoryControllerSuiteTest))
}
