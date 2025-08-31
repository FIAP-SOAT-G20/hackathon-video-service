package dto

import (
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
)

type CreateVideoInput struct {
	CustomerID uint64
}

type UpdateVideoInput struct {
	ID         uint64
	CustomerID uint64
	Status     valueobject.VideoStatus
	StaffID    uint64
}

type GetVideoInput struct {
	ID uint64
}

type DeleteVideoInput struct {
	ID uint64
}

type ListVideosInput struct {
	CustomerID    uint64
	Status        []valueobject.VideoStatus
	StatusExclude []valueobject.VideoStatus
	Page          int
	Limit         int
	Sort          string
}
