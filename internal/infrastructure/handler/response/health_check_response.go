package response

import "time"

type HealthCheckResponse struct {
	Status HealthCheckStatus                   `json:"status"`
	Checks map[string]HealthCheckVerifications `json:"checks,omitempty"`
}

type HealthCheckVerifications struct {
	ComponentId string            `json:"componentId"`
	Status      HealthCheckStatus `json:"status"`
	Time        time.Time         `json:"time"`
}

// HealthCheckStatus represents the status of a health check
type HealthCheckStatus string

// HealthCheckStatus values
const (
	HealthCheckStatusPass HealthCheckStatus = "pass"
	HealthCheckStatusWarn HealthCheckStatus = "warn"
	HealthCheckStatusFail HealthCheckStatus = "fail"
)

type HealthCheckLivenessResponse struct {
	Status string `json:"status"`
}
