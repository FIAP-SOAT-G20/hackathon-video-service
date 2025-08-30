package presenter

import (
	"encoding/json"
	"errors"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type orderJsonPresenter struct{}

// OrderJsonResponse represents the response of a order
func NewOrderJsonPresenter() port.Presenter {
	return &orderJsonPresenter{}
}

// Present write the response to the client
func (p *orderJsonPresenter) Present(pp dto.PresenterInput) ([]byte, error) {
	switch v := pp.Result.(type) {
	case *entity.Order:
		output := ToOrderJsonResponse(v)
		return json.Marshal(output)
	case []*entity.Order:
		orderOutputs := make([]OrderJsonResponse, len(v))
		for i, order := range v {
			orderOutputs[i] = ToOrderJsonResponse(order)
		}

		output := &OrderJsonPaginatedResponse{
			JsonPagination: JsonPagination{
				Total: pp.Total,
				Page:  pp.Page,
				Limit: pp.Limit,
			},
			Orders: orderOutputs,
		}
		return json.Marshal(output)
	default:
		return nil, domain.NewInternalError(errors.New(domain.ErrInternalError))
	}
}

// ToOrderJsonResponse convert entity.Order to OrderJsonResponse
func ToOrderJsonResponse(order *entity.Order) OrderJsonResponse {
	return OrderJsonResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Status:     string(order.Status),
		CreatedAt:  order.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  order.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}
