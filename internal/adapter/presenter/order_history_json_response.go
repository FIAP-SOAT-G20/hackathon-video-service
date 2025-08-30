package presenter

type OrderHistoryJsonResponse struct {
	ID        uint64  `json:"id" example:"1"`
	OrderID   uint64  `json:"order_id" example:"1"`
	StaffID   *uint64 `json:"staff_id" example:"1"`
	Status    string  `json:"status" example:"OPEN, CANCELLED, PENDING, RECEIVED, PREPARING, READY, COMPLETED"`
	CreatedAt string  `json:"created_at" example:"2024-02-09T10:00:00Z"`
}

type OrderHistoryJsonPaginatedResponse struct {
	JsonPagination
	OrderHistories []OrderHistoryJsonResponse `json:"order_histories"`
}
