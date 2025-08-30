package presenter

import (
	"encoding/json"
	"errors"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type orderHistoryJsonPresenter struct{}

// OrderHistoryJsonResponse represents the response of a orderHistory
func NewOrderHistoryJsonPresenter() port.Presenter {
	return &orderHistoryJsonPresenter{}
}

// toOrderHistoryJsonResponse convert entity.OrderHistory to OrderHistoryJsonResponse
func toOrderHistoryJsonResponse(orderHistory *entity.OrderHistory) OrderHistoryJsonResponse {
	return OrderHistoryJsonResponse{
		ID:        orderHistory.ID,
		OrderID:   orderHistory.OrderID,
		StaffID:   orderHistory.StaffID,
		Status:    orderHistory.Status.String(),
		CreatedAt: orderHistory.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Present write the response to the client
func (p *orderHistoryJsonPresenter) Present(pp dto.PresenterInput) ([]byte, error) {
	switch v := pp.Result.(type) {
	case *entity.OrderHistory:
		output := toOrderHistoryJsonResponse(v)
		return json.Marshal(output)
	case []*entity.OrderHistory:
		orderHistoryOutputs := make([]OrderHistoryJsonResponse, len(v))
		for i, orderHistory := range v {
			orderHistoryOutputs[i] = toOrderHistoryJsonResponse(orderHistory)
		}

		output := &OrderHistoryJsonPaginatedResponse{
			JsonPagination: JsonPagination{
				Total: pp.Total,
				Page:  pp.Page,
				Limit: pp.Limit,
			},
			OrderHistories: orderHistoryOutputs,
		}
		return json.Marshal(output)
	default:
		return nil, domain.NewInternalError(errors.New(domain.ErrInternalError))
	}
}
