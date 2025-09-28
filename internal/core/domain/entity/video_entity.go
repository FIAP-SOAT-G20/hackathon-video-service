package entity

import (
	"time"

	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
)

type Video struct {
	ID           uint64
	UserID       uint64
	Name         string
	Description  string
	Status       valueobject.VideoStatus
	Hash         string
	PresignedURL string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type VideoProcessedDownload struct {
	URL string
}

func (p *Video) Update(customerID uint64, status valueobject.VideoStatus, name string, description string, hash *string) {
	if customerID != 0 {
		p.UserID = customerID
	}
	if status != valueobject.UNDEFINED {
		p.Status = status
	}
	if name != "" {
		p.Name = name
	}
	if description != "" {
		p.Description = description
	}
	if hash != nil {
		p.Hash = *hash
	}
	p.UpdatedAt = time.Now()
}
