package presenter

import (
	"encoding/json"
	"errors"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type categoryJsonPresenter struct{}

// CategoryJsonResponse represents the response of a category
func NewCategoryJsonPresenter() port.Presenter {
	return &categoryJsonPresenter{}
}

// ToCategoryJsonResponse convert entity.Category to CategoryJsonResponse
func ToCategoryJsonResponse(category *entity.Category) CategoryJsonResponse {
	return CategoryJsonResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Present write the response to the client
func (p *categoryJsonPresenter) Present(pp dto.PresenterInput) ([]byte, error) {
	switch v := pp.Result.(type) {
	case *entity.Category:
		output := ToCategoryJsonResponse(v)
		return json.Marshal(output)
	case []*entity.Category:
		categoryOutputs := make([]CategoryJsonResponse, len(v))
		for i, category := range v {
			categoryOutputs[i] = ToCategoryJsonResponse(category)
		}

		output := &CategoryJsonPaginatedResponse{
			JsonPagination: JsonPagination{
				Total: pp.Total,
				Page:  pp.Page,
				Limit: pp.Limit,
			},
			Categories: categoryOutputs,
		}

		return json.Marshal(output)
	default:
		return nil, domain.NewInternalError(errors.New(domain.ErrInternalError))
	}
}
