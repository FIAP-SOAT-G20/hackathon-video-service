package presenter

import (
	"encoding/json"
	"errors"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type orderProductJsonPresenter struct{}

// OrderProductJsonResponse represents the response of a orderProduct
func NewOrderProductJsonPresenter() port.Presenter {
	return &orderProductJsonPresenter{}
}

// Present write the response to the client
func (p *orderProductJsonPresenter) Present(pp dto.PresenterInput) ([]byte, error) {
	switch v := pp.Result.(type) {
	case *entity.OrderProduct:
		output := ToOrderProductJsonResponse(v)
		return json.Marshal(output)
	case []*entity.OrderProduct:
		orderProductOutputs := make([]OrderProductJsonResponse, len(v))
		for i, orderProduct := range v {
			orderProductOutputs[i] = ToOrderProductJsonResponse(orderProduct)
		}

		output := &OrderProductJsonPaginatedResponse{
			JsonPagination: JsonPagination{
				Total: pp.Total,
				Page:  pp.Page,
				Limit: pp.Limit,
			},
			OrderProducts: orderProductOutputs,
		}
		return json.Marshal(output)
	default:
		return nil, domain.NewInternalError(errors.New(domain.ErrInternalError))
	}
}

// ToOrderProductJsonResponse convert entity.OrderProduct to OrderProductJsonResponse
func ToOrderProductJsonResponse(orderProduct *entity.OrderProduct) OrderProductJsonResponse {
	order := ToOrderJsonResponse(&orderProduct.Order)
	order.TotalBill = ""
	return OrderProductJsonResponse{
		OrderID:   orderProduct.OrderID,
		ProductID: orderProduct.ProductID,
		Quantity:  orderProduct.Quantity,
		Order:     order,
		Product:   ToProductJsonResponse(&orderProduct.Product),
		CreatedAt: orderProduct.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: orderProduct.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}
