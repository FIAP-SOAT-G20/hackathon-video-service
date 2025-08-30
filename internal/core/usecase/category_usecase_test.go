package usecase_test

import (
	"context"
	"testing"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (s *CategoryUsecaseSuiteTest) TestCategoryUseCase_List() {
	tests := []struct {
		name        string
		input       dto.ListCategoriesInput
		setupMocks  func()
		checkResult func(*testing.T, []*entity.Category, int64, error)
	}{
		{
			name: "should list categories successfully",
			input: dto.ListCategoriesInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, "", 1, 10).
					Return(s.mockCategories, int64(2), nil)
			},
			checkResult: func(t *testing.T, categories []*entity.Category, total int64, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, categories)
				assert.Equal(t, len(s.mockCategories), len(categories))
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should return error when repository fails",
			input: dto.ListCategoriesInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, "", 1, 10).
					Return(nil, int64(0), assert.AnError)
			},
			checkResult: func(t *testing.T, categories []*entity.Category, total int64, err error) {
				assert.Error(t, err)
				assert.Nil(t, categories)
				assert.Equal(t, int64(0), total)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name: "should filter by name",
			input: dto.ListCategoriesInput{
				Name:  "Test",
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, "Test", 1, 10).
					Return(s.mockCategories, int64(2), nil)
			},
			checkResult: func(t *testing.T, categories []*entity.Category, total int64, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, categories)
				assert.Equal(t, len(s.mockCategories), len(categories))
				assert.Equal(t, int64(2), total)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			categories, total, err := s.useCase.List(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, categories, total, err)
		})
	}
}

func (s *CategoryUsecaseSuiteTest) TestCategoryUseCase_Create() {
	tests := []struct {
		name        string
		input       dto.CreateCategoryInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Category, error)
	}{
		{
			name: "should create category successfully",
			input: dto.CreateCategoryInput{
				Name: s.mockCategories[0].Name,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(nil)
			},
			checkResult: func(t *testing.T, category *entity.Category, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, category)
				assert.Equal(t, s.mockCategories[0].Name, category.Name)
			},
		},
		{
			name: "should return error when gateway fails",
			input: dto.CreateCategoryInput{
				Name: s.mockCategories[0].Name,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, category *entity.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, category)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			category, err := s.useCase.Create(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, category, err)
		})
	}
}

func (s *CategoryUsecaseSuiteTest) TestCategoryUseCase_Get() {
	tests := []struct {
		name        string
		input       dto.GetCategoryInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Category, error)
	}{
		{
			name:  "should get category successfully",
			input: dto.GetCategoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockCategories[0], nil)
			},
			checkResult: func(t *testing.T, category *entity.Category, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, category)
				assert.Equal(t, s.mockCategories[0].ID, category.ID)
				assert.Equal(t, s.mockCategories[0].Name, category.Name)
				assert.Equal(t, s.mockCategories[0].CreatedAt, category.CreatedAt)
				assert.Equal(t, s.mockCategories[0].UpdatedAt, category.UpdatedAt)
			},
		},
		{
			name:  "should return not found error when category doesn't exist",
			input: dto.GetCategoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, categories *entity.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, categories)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name:  "should return internal error when gateway fails",
			input: dto.GetCategoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, category *entity.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, category)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			category, err := s.useCase.Get(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, category, err)
		})
	}
}

func (s *CategoryUsecaseSuiteTest) TestCategoryUseCase_Update() {
	tests := []struct {
		name        string
		input       dto.UpdateCategoryInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Category, error)
	}{
		{
			name: "should update category successfully",
			input: dto.UpdateCategoryInput{
				ID:   1,
				Name: "Foods UPDATED",
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockCategories[0], nil)

				s.mockGateway.EXPECT().
					Update(s.ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, p *entity.Category) error {
						assert.Equal(s.T(), "Foods UPDATED", p.Name)
						return nil
					})
			},
			checkResult: func(t *testing.T, category *entity.Category, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, category)
				assert.Equal(t, "Foods UPDATED", category.Name)
				assert.Equal(t, s.mockCategories[0].CreatedAt, category.CreatedAt)
			},
		},
		{
			name: "should return error when category not found",
			input: dto.UpdateCategoryInput{
				ID:   1,
				Name: "Foods UPDATED",
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, category *entity.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, category)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name: "should return error when gateway find fails",
			input: dto.UpdateCategoryInput{
				ID:   1,
				Name: "Foods UPDATED",
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, category *entity.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, category)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name: "should return error when gateway update fails",
			input: dto.UpdateCategoryInput{
				ID:   1,
				Name: "Foods UPDATED",
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockCategories[0], nil)

				s.mockGateway.EXPECT().
					Update(s.ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, category *entity.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, category)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			category, err := s.useCase.Update(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, category, err)
		})
	}
}

func (s *CategoryUsecaseSuiteTest) TestCategoryUseCase_Delete() {
	tests := []struct {
		name        string
		input       dto.DeleteCategoryInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Category, error)
	}{
		{
			name:  "should delete category successfully",
			input: dto.DeleteCategoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(&entity.Category{ID: 1}, nil)

				s.mockGateway.EXPECT().
					Delete(s.ctx, uint64(1)).
					Return(nil)
			},
			checkResult: func(t *testing.T, category *entity.Category, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, category)
				assert.Equal(t, uint64(1), category.ID)
			},
		},
		{
			name:  "should return not found error when category doesn't exist",
			input: dto.DeleteCategoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, category *entity.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, category)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name:  "should return error when gateway fails on find",
			input: dto.DeleteCategoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, category *entity.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, category)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name:  "should return error when gateway fails on delete",
			input: dto.DeleteCategoryInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(&entity.Category{}, nil)

				s.mockGateway.EXPECT().
					Delete(s.ctx, uint64(1)).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, category *entity.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, category)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			category, err := s.useCase.Delete(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, category, err)
		})
	}
}
