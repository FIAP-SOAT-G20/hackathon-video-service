package dto

import "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"

type GetCategoryInput struct {
	ID uint64
}

type ListCategoriesInput struct {
	Name  string
	Page  int
	Limit int
}

type UpdateCategoryInput struct {
	ID   uint64
	Name string
}

type DeleteCategoryInput struct {
	ID uint64
}

type CreateCategoryInput struct {
	Name string
}

func (c CreateCategoryInput) ToEntity() *entity.Category {
	return &entity.Category{
		Name: c.Name,
	}
}
