package controller_test

import (
	"testing"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/util"
	"github.com/stretchr/testify/assert"
)

func (s *CategoryControllerSuiteTest) TestCategoryController_ListCategories() {
	tests := []struct {
		name        string
		input       dto.ListCategoriesInput
		setupMocks  func()
		checkResult func(*testing.T, []byte, error)
	}{
		{
			name: "List categories success",
			input: dto.ListCategoriesInput{
				Name:  "Test",
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockUseCase.EXPECT().
					List(s.ctx, dto.ListCategoriesInput{
						Name:  "Test",
						Page:  1,
						Limit: 10,
					}).
					Return(s.mockCategories, int64(2), nil)
			},
			checkResult: func(t *testing.T, output []byte, err error) {
				want, _ := util.ReadGoldenFile("category/list_success")
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, want, util.RemoveAllSpaces(string(output)))
			},
		},
		{
			name: "List categories use case error",
			input: dto.ListCategoriesInput{
				Name:  "Test",
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockUseCase.EXPECT().
					List(s.ctx, dto.ListCategoriesInput{
						Name:  "Test",
						Page:  1,
						Limit: 10,
					}).
					Return(nil, int64(0), assert.AnError)
			},
			checkResult: func(t *testing.T, output []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			output, err := s.controller.List(s.ctx, s.mockPresenter, tt.input)

			// Assert
			tt.checkResult(t, output, err)
		})
	}
}

func (s *CategoryControllerSuiteTest) TestCategoryController_GetCategory() {
	tests := []struct {
		name        string
		input       dto.GetCategoryInput
		setupMocks  func()
		checkResult func(*testing.T, []byte, error)
	}{
		{
			name:  "Get category success",
			input: dto.GetCategoryInput{ID: uint64(1)},
			setupMocks: func() {
				s.mockUseCase.EXPECT().
					Get(s.ctx, dto.GetCategoryInput{
						ID: uint64(1),
					}).
					Return(s.mockCategory, nil)
			},
			checkResult: func(t *testing.T, output []byte, err error) {
				want, _ := util.ReadGoldenFile("category/get_success")
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, want, util.RemoveAllSpaces(string(output)))
			},
		},
		{
			name:  "Get use case error",
			input: dto.GetCategoryInput{ID: uint64(1)},
			setupMocks: func() {
				s.mockUseCase.EXPECT().
					Get(s.ctx, dto.GetCategoryInput{ID: uint64(1)}).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, output []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			output, err := s.controller.Get(s.ctx, s.mockPresenter, tt.input)

			// Assert
			tt.checkResult(t, output, err)
		})
	}
}

func (s *CategoryControllerSuiteTest) TestCategoryController_UpdateCategory() {
	categoryUpdated := &entity.Category{
		ID:        6,
		Name:      "Foods UPDATED",
		CreatedAt: s.mockCategory.CreatedAt,
		UpdatedAt: s.mockCategory.UpdatedAt,
	}
	tests := []struct {
		name        string
		input       dto.UpdateCategoryInput
		setupMocks  func()
		checkResult func(*testing.T, []byte, error)
	}{
		{
			name: "Update category success",
			input: dto.UpdateCategoryInput{
				ID:   categoryUpdated.ID,
				Name: categoryUpdated.Name,
			},
			setupMocks: func() {
				s.mockUseCase.EXPECT().
					Update(s.ctx, dto.UpdateCategoryInput{
						ID:   categoryUpdated.ID,
						Name: categoryUpdated.Name,
					}).
					Return(categoryUpdated, nil)
			},
			checkResult: func(t *testing.T, output []byte, err error) {
				want, _ := util.ReadGoldenFile("category/update_success")
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, want, util.RemoveAllSpaces(string(output)))
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			output, err := s.controller.Update(s.ctx, s.mockPresenter, tt.input)

			// Assert
			tt.checkResult(t, output, err)
		})
	}
}

func (s *CategoryControllerSuiteTest) TestCategoryController_DeleteCategory() {
	tests := []struct {
		name        string
		input       dto.DeleteCategoryInput
		setupMocks  func()
		checkResult func(*testing.T, []byte, error)
	}{
		{
			name:  "Delete category success",
			input: dto.DeleteCategoryInput{ID: uint64(1)},
			setupMocks: func() {
				s.mockUseCase.EXPECT().
					Delete(s.ctx, dto.DeleteCategoryInput{ID: uint64(1)}).
					Return(s.mockCategory, nil)
			},
			checkResult: func(t *testing.T, output []byte, err error) {
				want, _ := util.ReadGoldenFile("category/delete_success")
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, want, util.RemoveAllSpaces(string(output)))
			},
		},
		{
			name:  "Delete use case error",
			input: dto.DeleteCategoryInput{ID: uint64(1)},
			setupMocks: func() {
				s.mockUseCase.EXPECT().
					Delete(s.ctx, dto.DeleteCategoryInput{ID: uint64(1)}).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, output []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			output, err := s.controller.Delete(s.ctx, s.mockPresenter, tt.input)

			// Assert
			tt.checkResult(t, output, err)
		})
	}
}

func (s *CategoryControllerSuiteTest) TestCategoryController_CreateCategory() {
	tests := []struct {
		name        string
		input       dto.CreateCategoryInput
		setupMocks  func()
		checkResult func(*testing.T, []byte, error)
	}{
		{
			name: "Create category success",
			input: dto.CreateCategoryInput{
				Name: s.mockCategory.Name,
			},
			setupMocks: func() {
				s.mockUseCase.EXPECT().
					Create(s.ctx, dto.CreateCategoryInput{
						Name: s.mockCategory.Name,
					}).
					Return(s.mockCategory, nil)
			},
			checkResult: func(t *testing.T, output []byte, err error) {
				want, _ := util.ReadGoldenFile("category/create_success")
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, want, util.RemoveAllSpaces(string(output)))
			},
		},
		{
			name: "Create use case error",
			input: dto.CreateCategoryInput{
				Name: s.mockCategory.Name,
			},
			setupMocks: func() {
				s.mockUseCase.EXPECT().
					Create(s.ctx, dto.CreateCategoryInput{
						Name: s.mockCategory.Name,
					}).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, output []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			output, err := s.controller.Create(s.ctx, s.mockPresenter, tt.input)

			// Assert
			tt.checkResult(t, output, err)
		})
	}
}
