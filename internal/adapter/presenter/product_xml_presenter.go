package presenter

import (
	"encoding/xml"
	"errors"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type productXmlPresenter struct{}

// NewProductXmlPresenter creates a new ProductXmlPresenter
func NewProductXmlPresenter() port.Presenter {
	return &productXmlPresenter{}
}

// Present writes the response to the client
func (p *productXmlPresenter) Present(pp dto.PresenterInput) ([]byte, error) {
	switch v := pp.Result.(type) {
	case *entity.Product:
		output := toProductXmlResponse(v)
		return xml.Marshal(output)
	case []*entity.Product:
		productOutputs := make([]ProductXmlResponse, len(v))
		for i, product := range v {
			productOutputs[i] = toProductXmlResponse(product)
		}

		output := &ProductXmlPaginatedResponse{
			XmlPagination: XmlPagination{
				Total: pp.Total,
				Page:  pp.Page,
				Limit: pp.Limit,
			},
			Products: productOutputs,
		}
		return xml.Marshal(output)
	default:
		return nil, domain.NewInternalError(errors.New(domain.ErrInternalError))
	}
}

// toProductXmlResponse converts a Product entity to a ProductXmlResponse
func toProductXmlResponse(product *entity.Product) ProductXmlResponse {
	return ProductXmlResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}
