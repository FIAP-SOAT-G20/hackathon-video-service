package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

type CategoryUseCase interface {
	Create(ctx context.Context, input dto.CreateCategoryInput) (*entity.Category, error)
	Get(ctx context.Context, input dto.GetCategoryInput) (*entity.Category, error)
	List(ctx context.Context, input dto.ListCategoriesInput) ([]*entity.Category, int64, error)
	Update(ctx context.Context, input dto.UpdateCategoryInput) (*entity.Category, error)
	Delete(ctx context.Context, input dto.DeleteCategoryInput) (*entity.Category, error)
}
