package entity

import valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"

type VideoStatusUpdated struct {
	VideoID uint64                  `json:"video_id"`
	Status  valueobject.VideoStatus `json:"status"`
	StaffID *uint64                 `json:"staff_id,omitempty"` // Optional field for staff ID
}
