package valueobject

import (
	"slices"
	"strings"
)

type VideoStatus string

const (
	OPEN       VideoStatus = "OPEN"
	CANCELLED  VideoStatus = "CANCELLED"
	PENDING    VideoStatus = "PENDING"
	RECEIVED   VideoStatus = "RECEIVED"
	PREPARING  VideoStatus = "PREPARING"
	READY      VideoStatus = "READY"
	COMPLETED  VideoStatus = "COMPLETED"
	UNDEFINDED VideoStatus = "UNDEFINDED"
)

func IsValidVideoStatus(status string) bool {
	_, ok := ToVideoStatus(status)
	return ok
}

// String returns the string representation of the VideoStatus
func (o VideoStatus) String() string {
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

// ToVideoStatus converts a string to a VideoStatus
func ToVideoStatus(status string) (VideoStatus, bool) {
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

// VideoStatusTransitions defines the allowed transitions between VideoStatuses
var VideoStatusTransitions = map[VideoStatus][]VideoStatus{
	OPEN:      {CANCELLED, PENDING, RECEIVED},
	CANCELLED: {},
	PENDING:   {OPEN, RECEIVED, CANCELLED},
	RECEIVED:  {PREPARING, CANCELLED},
	PREPARING: {READY, CANCELLED},
	READY:     {COMPLETED},
	COMPLETED: {},
}

// StatusCanTransitionTo returns true if the transition from oldStatus to newStatus is allowed
func StatusCanTransitionTo(oldStatus, newStatus VideoStatus) bool {
	allowedStatuses := VideoStatusTransitions[oldStatus]
	return slices.Contains(allowedStatuses, newStatus)
}

// StatusTransitionNeedsStaffID returns true if the new status requires a staff ID
func StatusTransitionNeedsStaffID(newStatus VideoStatus) bool {
	switch newStatus {
	case OPEN, CANCELLED, PENDING, RECEIVED:
		return false
	case PREPARING, READY, COMPLETED:
		return true
	default:
		return false
	}
}
