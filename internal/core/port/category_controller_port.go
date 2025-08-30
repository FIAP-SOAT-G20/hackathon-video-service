package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

type CategoryController interface {
	Create(ctx context.Context, presenter Presenter, input dto.CreateCategoryInput) ([]byte, error)
	List(ctx context.Context, presenter Presenter, input dto.ListCategoriesInput) ([]byte, error)
	Get(ctx context.Context, presenter Presenter, input dto.GetCategoryInput) ([]byte, error)
	Update(ctx context.Context, presenter Presenter, input dto.UpdateCategoryInput) ([]byte, error)
	Delete(ctx context.Context, presenter Presenter, input dto.DeleteCategoryInput) ([]byte, error)
}
