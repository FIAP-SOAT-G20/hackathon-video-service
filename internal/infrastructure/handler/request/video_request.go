package request

import valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"

type ListVideosQueryRequest struct {
	UserID        uint64 `form:"user_id" example:"1" default:"0"`
	Status        string `form:"status" binding:"omitempty" example:"PROCESSING"`
	StatusExclude string `form:"status_exclude" binding:"omitempty" example:"FAILED,FINISHED"`
	Page          int    `form:"page,default=1" example:"1"`
	Limit         int    `form:"limit,default=10" example:"10"`
	// Sort by default: status:d,created_at. Use <field_name>:d for descending, and the default order is ascending
	Sort string `form:"sort" example:"status:d,created_at"`
}

type CreateVideoBodyRequest struct {
	UserID      uint64 `json:"user_id" binding:"required" example:"1"`
	Name        string `json:"name" binding:"required" example:"My Video"`
	Description string `json:"description" binding:"omitempty" example:"This is my video description"`
}

type GetVideoUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type UpdateVideoUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type UpdateVideoBodyRequest struct {
	Status valueobject.VideoStatus `json:"status" binding:"required,video_status_exists" example:"PROCESSING"`
	Hash   string                  `json:"hash" binding:"omitempty" example:"abc123hash456"`
}

type UpdateVideoPartilRequest struct {
	Status valueobject.VideoStatus `json:"status" example:"PROCESSING"`
}

type UpdateVideoPartilBodyRequest struct {
	Status valueobject.VideoStatus `json:"status" binding:"omitempty,video_status_exists" example:"PROCESSING"`
	Hash   string                  `json:"hash" binding:"omitempty" example:"abc123hash456"`
}

type DeleteVideoUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type DownloadVideoUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}
