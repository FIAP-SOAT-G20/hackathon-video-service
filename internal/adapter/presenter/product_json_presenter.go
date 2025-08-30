package presenter

import (
	"encoding/json"
	"errors"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type productJsonPresenter struct{}

// ProductJsonResponse represents the response of a product
func NewProductJsonPresenter() port.Presenter {
	return &productJsonPresenter{}
}

// Present write the response to the client
func (p *productJsonPresenter) Present(pp dto.PresenterInput) ([]byte, error) {
	switch v := pp.Result.(type) {
	case *entity.Product:
		output := ToProductJsonResponse(v)
		return json.Marshal(output)
	case []*entity.Product:
		productOutputs := make([]ProductJsonResponse, len(v))
		for i, product := range v {
			productOutputs[i] = ToProductJsonResponse(product)
		}

		output := &ProductJsonPaginatedResponse{
			JsonPagination: JsonPagination{
				Total: pp.Total,
				Page:  pp.Page,
				Limit: pp.Limit,
			},
			Products: productOutputs,
		}
		return json.Marshal(output)
	default:
		return nil, domain.NewInternalError(errors.New(domain.ErrInternalError))
	}
}

// ToProductJsonResponse convert entity.Product to ProductJsonResponse
func ToProductJsonResponse(product *entity.Product) ProductJsonResponse {
	return ProductJsonResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}
