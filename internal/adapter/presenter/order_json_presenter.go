package presenter

import (
	"encoding/json"
	"errors"
	"fmt"

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
		TotalBill:  calculateTotalBill(order.OrderProducts),
		Status:     string(order.Status),
		Products:   ToProductsJsonResponse(order.OrderProducts),
		CreatedAt:  order.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  order.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToProductsJsonResponse convert a slice of entity.OrderProduct to a slice of ProductsJsonResponse
func ToProductsJsonResponse(orderProducts []entity.OrderProduct) []ProductsJsonResponse {
	products := make([]ProductsJsonResponse, len(orderProducts))
	for i, orderProduct := range orderProducts {
		products[i] = ProductsJsonResponse{
			ProductJsonResponse: ToProductJsonResponse(&orderProduct.Product),
			Quantity:            orderProduct.Quantity,
		}
	}
	return products
}

// calculateTotalBill calculate the total bill of an order
func calculateTotalBill(orderProducts []entity.OrderProduct) string {
	var total float64
	for _, orderProduct := range orderProducts {
		total += orderProduct.Product.Price * float64(orderProduct.Quantity)
	}
	// 2 decimal places
	return fmt.Sprintf("%.2f", total)
}
