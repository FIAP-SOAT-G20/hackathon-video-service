package dto

import (
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
)

type CreateOrderInput struct {
	CustomerID uint64
}

type UpdateOrderInput struct {
	ID         uint64
	CustomerID uint64
	Status     valueobject.OrderStatus
	StaffID    uint64
}

type GetOrderInput struct {
	ID uint64
}

type DeleteOrderInput struct {
	ID uint64
}

type ListOrdersInput struct {
	CustomerID    uint64
	Status        []valueobject.OrderStatus
	StatusExclude []valueobject.OrderStatus
	Page          int
	Limit         int
	Sort          string
}
