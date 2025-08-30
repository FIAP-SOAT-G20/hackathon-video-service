package request

import valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"

type ListOrdersQueryRequest struct {
	CustomerID    uint64 `form:"customer_id" example:"1" default:"0"`
	Status        string `form:"status" binding:"omitempty" example:"PENDING"`
	StatusExclude string `form:"status_exclude" binding:"omitempty" example:"CANCELLED,COMPLETED"`
	Page          int    `form:"page,default=1" example:"1"`
	Limit         int    `form:"limit,default=10" example:"10"`
	// Sort by default: status:d,created_at. Use <field_name>:d for descending, and the default order is ascending
	Sort string `form:"sort" example:"status:d,created_at"`
}

type CreateOrderBodyRequest struct {
	CustomerID uint64 `json:"customer_id" binding:"required" example:"1"`
}

type GetOrderUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type UpdateOrderUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type UpdateOrderBodyRequest struct {
	// StaffID is only required when status is PREPARING, READY or COMPLETED
	StaffID    uint64                  `json:"staff_id" example:"1"`
	CustomerID uint64                  `json:"customer_id" binding:"required" example:"1"`
	Status     valueobject.OrderStatus `json:"status" binding:"required,order_status_exists" example:"PENDING"`
}

type UpdateOrderPartilRequest struct {
	// StaffID is only required when status is PREPARING, READY or COMPLETED
	StaffID uint64                  `json:"staff_id" example:"1"`
	Status  valueobject.OrderStatus `json:"status" example:"PENDING"`
}

type UpdateOrderPartilBodyRequest struct {
	CustomerID uint64 `json:"customer_id" example:"1"`
	// StaffID is only required when status is PREPARING, READY or COMPLETED
	StaffID uint64                  `json:"staff_id" example:"1"`
	Status  valueobject.OrderStatus `json:"status" binding:"omitempty,order_status_exists" example:"PENDING"`
}

type DeleteOrderUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}
