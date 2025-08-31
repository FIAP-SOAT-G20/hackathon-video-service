package entity

import (
	"time"

	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
)

type Video struct {
	ID         uint64
	CustomerID uint64
	Status     valueobject.VideoStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (p *Video) Update(customerID uint64, status valueobject.VideoStatus) {
	if customerID != 0 {
		p.CustomerID = customerID
	}
	if status != valueobject.UNDEFINDED {
		p.Status = status
	}
	p.UpdatedAt = time.Now()
}
