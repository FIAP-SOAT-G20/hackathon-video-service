package entity

import valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"

type VideoStatusUpdated struct {
	ID     uint64                  `json:"id"`
	Status valueobject.VideoStatus `json:"status"`
}
