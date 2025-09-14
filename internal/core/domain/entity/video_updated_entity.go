package entity

import valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"

type VideoUpdated struct {
	VideoID uint64                  `json:"video_id"`
	Status  valueobject.VideoStatus `json:"status"`
	Hash    string                  `json:"hash"`
}
