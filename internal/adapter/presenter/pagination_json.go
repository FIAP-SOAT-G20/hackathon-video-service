package presenter

type JsonPagination struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	TotalPages int64 `json:"total_pages" example:"100"`
	TotalItems int64 `json:"total_items" example:"1000"`
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
}
