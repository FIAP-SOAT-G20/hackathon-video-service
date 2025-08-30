package dto

import valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"

type CreateFirstOrderHistoryInput struct {
	OrderID uint64
}

type CreateOrderHistoryInput struct {
	OrderID uint64
	StaffID *uint64
	Status  valueobject.OrderStatus
}

type GetOrderHistoryInput struct {
	ID uint64
}

type DeleteOrderHistoryInput struct {
	ID uint64
}

type ListOrderHistoriesInput struct {
	OrderID uint64
	Status  valueobject.OrderStatus
	Page    int
	Limit   int
}
