package valueobject

import (
	"slices"
	"strings"
)

type VideoStatus string

const (
	CREATED    VideoStatus = "CREATED"
	UPLOADED   VideoStatus = "UPLOADED"
	PROCESSING VideoStatus = "PROCESSING"
	FINISHED   VideoStatus = "FINISHED"
	FAILED     VideoStatus = "FAILED"
	UNDEFINED  VideoStatus = "UNDEFINED"
)

func IsValidVideoStatus(status string) bool {
	_, ok := ToVideoStatus(status)
	return ok
}

// String returns the string representation of the VideoStatus
func (o VideoStatus) String() string {
	switch o {
	case CREATED:
		return "CREATED"
	case UPLOADED:
		return "UPLOADED"
	case PROCESSING:
		return "PROCESSING"
	case FINISHED:
		return "FINISHED"
	case FAILED:
		return "FAILED"
	default:
		return "UNDEFINED"
	}
}

// ToVideoStatus converts a string to a VideoStatus
func ToVideoStatus(status string) (VideoStatus, bool) {
	switch strings.ToUpper(status) {
	case "CREATED":
		return CREATED, true
	case "UPLOADED":
		return UPLOADED, true
	case "PROCESSING":
		return PROCESSING, true
	case "FINISHED":
		return FINISHED, true
	case "FAILED":
		return FAILED, true
	default:
		return UNDEFINED, false
	}
}

// VideoStatusTransitions defines the allowed transitions between VideoStatuses
var VideoStatusTransitions = map[VideoStatus][]VideoStatus{
	CREATED:    {FAILED, UPLOADED},
	UPLOADED:   {FAILED, PROCESSING},
	PROCESSING: {FAILED, FINISHED},
	FINISHED:   {},
	FAILED:     {},
}

// StatusCanTransitionTo returns true if the transition from oldStatus to newStatus is allowed
func StatusCanTransitionTo(oldStatus, newStatus VideoStatus) bool {
	allowedStatuses := VideoStatusTransitions[oldStatus]
	return slices.Contains(allowedStatuses, newStatus)
}
