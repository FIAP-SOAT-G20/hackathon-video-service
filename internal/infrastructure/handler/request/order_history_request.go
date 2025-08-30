package request

import valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"

type ListOrderHistoriesQueryRequest struct {
	OrderID uint64                  `form:"order_id,default=0" example:"1"`
	Status  valueobject.OrderStatus `form:"status" binding:"omitempty" example:"OPEN, CANCELLED, PENDING, RECEIVED, PREPARING, READY, COMPLETED"`
	Page    int                     `form:"page,default=1" example:"1"`
	Limit   int                     `form:"limit,default=10" example:"10"`
}

type GetOrderHistoryUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type DeleteOrderHistoryUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}
