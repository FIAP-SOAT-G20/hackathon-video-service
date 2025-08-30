package presenter

type OrderJsonResponse struct {
	ID         uint64                 `json:"id"`
	CustomerID uint64                 `json:"customer_id" example:"1"`
	TotalBill  string                 `json:"total_bill,omitempty" example:"100.00"`
	Status     string                 `json:"status" example:"PENDING"`
	Products   []ProductsJsonResponse `json:"products,omitempty"`
	CreatedAt  string                 `json:"created_at" example:"2024-02-09T10:00:00Z"`
	UpdatedAt  string                 `json:"updated_at" example:"2024-02-09T10:00:00Z"`
}

type OrderJsonPaginatedResponse struct {
	JsonPagination
	Orders []OrderJsonResponse `json:"orders"`
}

type ProductsJsonResponse struct {
	ProductJsonResponse
	Quantity uint32 `json:"quantity"`
}
