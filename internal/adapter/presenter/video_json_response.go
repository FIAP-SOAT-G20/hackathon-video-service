package presenter

type VideoJsonResponse struct {
	ID          uint64 `json:"id"`
	UserID      uint64 `json:"user_id" example:"1"`
	Name        string `json:"name" example:"My Video"`
	Description string `json:"description" example:"This is a description"`
	Status      string `json:"status" example:"PENDING"`
	Hash        string `json:"hash,omitempty" example:"hash_value"`
	Link        string `json:"link,omitempty" example:"https://example.com/video/1"`
	CreatedAt   string `json:"created_at" example:"2024-02-09T10:00:00Z"`
	UpdatedAt   string `json:"updated_at" example:"2024-02-09T10:00:00Z"`
}

type VideoJsonPaginatedResponse struct {
	JsonPagination
	Videos []VideoJsonResponse `json:"videos"`
}
