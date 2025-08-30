package valueobject

import (
	"slices"
	"strings"
)

type OrderStatus string

const (
	OPEN       OrderStatus = "OPEN"
	CANCELLED  OrderStatus = "CANCELLED"
	PENDING    OrderStatus = "PENDING"
	RECEIVED   OrderStatus = "RECEIVED"
	PREPARING  OrderStatus = "PREPARING"
	READY      OrderStatus = "READY"
	COMPLETED  OrderStatus = "COMPLETED"
	UNDEFINDED OrderStatus = "UNDEFINDED"
)

func IsValidOrderStatus(status string) bool {
	_, ok := ToOrderStatus(status)
	return ok
}

// String returns the string representation of the OrderStatus
func (o OrderStatus) String() string {
	switch o {
	case OPEN:
		return "OPEN"
	case CANCELLED:
		return "CANCELLED"
	case PENDING:
		return "PENDING"
	case RECEIVED:
		return "RECEIVED"
	case PREPARING:
		return "PREPARING"
	case READY:
		return "READY"
	case COMPLETED:
		return "COMPLETED"
	default:
		return "UNDEFINDED"
	}
}

// ToOrderStatus converts a string to an OrderStatus
func ToOrderStatus(status string) (OrderStatus, bool) {
	switch strings.ToUpper(status) {
	case "OPEN":
		return OPEN, true
	case "CANCELLED":
		return CANCELLED, true
	case "PENDING":
		return PENDING, true
	case "RECEIVED":
		return RECEIVED, true
	case "PREPARING":
		return PREPARING, true
	case "READY":
		return READY, true
	case "COMPLETED":
		return COMPLETED, true
	default:
		return UNDEFINDED, false
	}
}

// OrderStatusTransitions defines the allowed transitions between OrderStatuses
var OrderStatusTransitions = map[OrderStatus][]OrderStatus{
	OPEN:      {CANCELLED, PENDING, RECEIVED},
	CANCELLED: {},
	PENDING:   {OPEN, RECEIVED, CANCELLED},
	RECEIVED:  {PREPARING, CANCELLED},
	PREPARING: {READY, CANCELLED},
	READY:     {COMPLETED},
	COMPLETED: {},
}

// StatusCanTransitionTo returns true if the transition from oldStatus to newStatus is allowed
func StatusCanTransitionTo(oldStatus, newStatus OrderStatus) bool {
	allowedStatuses := OrderStatusTransitions[oldStatus]
	return slices.Contains(allowedStatuses, newStatus)
}

// StatusTransitionNeedsStaffID returns true if the new status requires a staff ID
func StatusTransitionNeedsStaffID(newStatus OrderStatus) bool {
	switch newStatus {
	case OPEN, CANCELLED, PENDING, RECEIVED:
		return false
	case PREPARING, READY, COMPLETED:
		return true
	default:
		return false
	}
}
