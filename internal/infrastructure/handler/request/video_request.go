package request

import valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"

type ListVideosQueryRequest struct {
	CustomerID    uint64 `form:"user_id" example:"1" default:"0"`
	Status        string `form:"status" binding:"omitempty" example:"PENDING"`
	StatusExclude string `form:"status_exclude" binding:"omitempty" example:"CANCELLED,COMPLETED"`
	Page          int    `form:"page,default=1" example:"1"`
	Limit         int    `form:"limit,default=10" example:"10"`
	// Sort by default: status:d,created_at. Use <field_name>:d for descending, and the default order is ascending
	Sort string `form:"sort" example:"status:d,created_at"`
}

type CreateVideoBodyRequest struct {
	CustomerID uint64 `json:"user_id" binding:"required" example:"1"`
}

type GetVideoUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type UpdateVideoUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type UpdateVideoBodyRequest struct {
	// StaffID is only required when status is PREPARING, READY or COMPLETED
	StaffID    uint64                  `json:"staff_id" example:"1"`
	CustomerID uint64                  `json:"user_id" binding:"required" example:"1"`
	Status     valueobject.VideoStatus `json:"status" binding:"required,video_status_exists" example:"PENDING"`
}

type UpdateVideoPartilRequest struct {
	// StaffID is only required when status is PREPARING, READY or COMPLETED
	StaffID uint64                  `json:"staff_id" example:"1"`
	Status  valueobject.VideoStatus `json:"status" example:"PENDING"`
}

type UpdateVideoPartilBodyRequest struct {
	CustomerID uint64 `json:"user_id" example:"1"`
	// StaffID is only required when status is PREPARING, READY or COMPLETED
	StaffID uint64                  `json:"staff_id" example:"1"`
	Status  valueobject.VideoStatus `json:"status" binding:"omitempty,video_status_exists" example:"PENDING"`
}

type DeleteVideoUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}
