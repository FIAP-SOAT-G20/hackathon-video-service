package entity

import valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"

type OrderStatusUpdated struct {
	OrderID uint64                  `json:"order_id"`
	Status  valueobject.OrderStatus `json:"status"`
	StaffID *uint64                 `json:"staff_id,omitempty"` // Optional field for staff ID
}
