package dto

import "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"

type CreateProductInput struct {
	Name        string
	Description string
	Price       float64
	CategoryID  uint64
}

func (i CreateProductInput) ToEntity() *entity.Product {
	return &entity.Product{
		Name:        i.Name,
		Description: i.Description,
		Price:       i.Price,
		CategoryID:  i.CategoryID,
	}
}

type UpdateProductInput struct {
	ID          uint64
	Name        string
	Description string
	Price       float64
	CategoryID  uint64
}

type GetProductInput struct {
	ID uint64
}

type DeleteProductInput struct {
	ID uint64
}

type ListProductsInput struct {
	Name       string
	CategoryID uint64
	Page       int
	Limit      int
}
