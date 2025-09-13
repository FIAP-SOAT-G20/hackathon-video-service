package dto

import (
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
)

type CreateVideoInput struct {
	UserID      uint64
	Name        string
	Description string
}

type UpdateVideoInput struct {
	ID          uint64
	UserID      uint64
	Name        string
	Description string
	Status      valueobject.VideoStatus
	StaffID     uint64
}

type GetVideoInput struct {
	ID uint64
}

type DeleteVideoInput struct {
	ID uint64
}

type ListVideosInput struct {
	UserID        uint64
	Status        []valueobject.VideoStatus
	StatusExclude []valueobject.VideoStatus
	Page          int
	Limit         int
	Sort          string
}

type DownloadVideoInput struct {
	ID uint64
}
