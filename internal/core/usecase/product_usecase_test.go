package usecase_test

import (
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

func (s *ProductUsecaseSuiteTest) TestProductsUseCase_List() {
	tests := []struct {
		name        string
		input       dto.ListProductsInput
		setupMocks  func()
		checkResult func(*testing.T, []*entity.Product, int64, error)
	}{
		{
			name: "should list products successfully",
			input: dto.ListProductsInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, "", uint64(0), 1, 10).
					Return(s.mockProducts, int64(2), nil)
			},
			checkResult: func(t *testing.T, products []*entity.Product, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockProducts, products)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should return error when repository fails",
			input: dto.ListProductsInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, "", uint64(0), 1, 10).
					Return(nil, int64(0), assert.AnError)
			},
			checkResult: func(t *testing.T, products []*entity.Product, total int64, err error) {
				assert.Error(t, err)
				assert.Nil(t, products)
				assert.Equal(t, int64(0), total)
			},
		},
		{
			name: "should filter by name",
			input: dto.ListProductsInput{
				Name:  "Test",
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, "Test", uint64(0), 1, 10).
					Return(s.mockProducts, int64(2), nil)
			},
			checkResult: func(t *testing.T, products []*entity.Product, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockProducts, products)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should filter by category",
			input: dto.ListProductsInput{
				CategoryID: 1,
				Page:       1,
				Limit:      10,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindAll(s.ctx, "", uint64(1), 1, 10).
					Return(s.mockProducts, int64(2), nil)
			},
			checkResult: func(t *testing.T, products []*entity.Product, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, s.mockProducts, products)
				assert.Equal(t, int64(2), total)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			products, total, err := s.useCase.List(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, products, total, err)
		})
	}
}

func (s *ProductUsecaseSuiteTest) TestProductUseCase_Create() {
	tests := []struct {
		name        string
		input       dto.CreateProductInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Product, error)
	}{
		{
			name: "should create product successfully",
			input: dto.CreateProductInput{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       99.99,
				CategoryID:  1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(nil)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, product)
				assert.Equal(t, "Test Product", product.Name)
				assert.Equal(t, "Test Description", product.Description)
				assert.Equal(t, 99.99, product.Price)
				assert.Equal(t, uint64(1), product.CategoryID)
			},
		},
		{
			name: "should return error when gateway fails",
			input: dto.CreateProductInput{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       99.99,
				CategoryID:  1,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.Error(t, err)
				assert.Nil(t, product)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			product, err := s.useCase.Create(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, product, err)
		})
	}
}

func (s *ProductUsecaseSuiteTest) TestProductUseCase_Get() {
	tests := []struct {
		name        string
		input       dto.GetProductInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Product, error)
	}{
		{
			name:  "should get product successfully",
			input: dto.GetProductInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockProducts[0], nil)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, product)
				assert.Equal(t, uint64(1), product.ID)
			},
		},
		{
			name:  "should return not found error when product doesn't exist",
			input: dto.GetProductInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.Error(t, err)
				assert.Nil(t, product)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name:  "should return internal error when gateway fails",
			input: dto.GetProductInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.Error(t, err)
				assert.Nil(t, product)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			product, err := s.useCase.Get(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, product, err)
		})
	}
}

func (s *ProductUsecaseSuiteTest) TestProductUseCase_Update() {
	tests := []struct {
		name        string
		input       dto.UpdateProductInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Product, error)
	}{
		{
			name: "should update product successfully",
			input: dto.UpdateProductInput{
				ID:          1,
				Name:        "New Name",
				Description: "New Description",
				Price:       20.0,
				CategoryID:  2,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockProducts[0], nil)

				s.mockGateway.EXPECT().
					Update(s.ctx, gomock.Any()).
					Return(nil)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, product)
				assert.Equal(t, "New Name", product.Name)
				assert.Equal(t, "New Description", product.Description)
				assert.Equal(t, 20.0, product.Price)
				assert.Equal(t, uint64(2), product.CategoryID)
			},
		},
		{
			name: "should return error when product not found",
			input: dto.UpdateProductInput{
				ID:          1,
				Name:        "New Name",
				Description: "New Description",
				Price:       20.0,
				CategoryID:  2,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.Error(t, err)
				assert.Nil(t, product)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name: "should return error when gateway find fails",
			input: dto.UpdateProductInput{
				ID:          1,
				Name:        "New Name",
				Description: "New Description",
				Price:       20.0,
				CategoryID:  2,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.Error(t, err)
				assert.Nil(t, product)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name: "should return error when gateway update fails",
			input: dto.UpdateProductInput{
				ID:          1,
				Name:        "New Name",
				Description: "New Description",
				Price:       20.0,
				CategoryID:  2,
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(s.mockProducts[0], nil)

				s.mockGateway.EXPECT().
					Update(s.ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.Error(t, err)
				assert.Nil(t, product)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			product, err := s.useCase.Update(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, product, err)
		})
	}
}

func (s *ProductUsecaseSuiteTest) TestProductUseCase_Delete() {
	tests := []struct {
		name        string
		input       dto.DeleteProductInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Product, error)
	}{
		{
			name:  "should delete _test successfully",
			input: dto.DeleteProductInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(&entity.Product{}, nil)

				s.mockGateway.EXPECT().
					Delete(s.ctx, uint64(1)).
					Return(nil)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, product)
			},
		},
		{
			name:  "should return not found error when _test doesn't exist",
			input: dto.DeleteProductInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.Error(t, err)
				assert.Nil(t, product)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name:  "should return error when gateway fails on find",
			input: dto.DeleteProductInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.Error(t, err)
				assert.Nil(t, product)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name:  "should return error when gateway fails on delete",
			input: dto.DeleteProductInput{ID: 1},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					FindByID(s.ctx, uint64(1)).
					Return(&entity.Product{}, nil)

				s.mockGateway.EXPECT().
					Delete(s.ctx, uint64(1)).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, product *entity.Product, err error) {
				assert.Error(t, err)
				assert.Nil(t, product)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			product, err := s.useCase.Delete(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, product, err)
		})
	}
}
