package presenter

type OrderProductJsonResponse struct {
	OrderID   uint64              `json:"order_id"`
	ProductID uint64              `json:"product_id"`
	Quantity  uint32              `json:"quantity"`
	Order     OrderJsonResponse   `json:"order,omitempty"`
	Product   ProductJsonResponse `json:"product,omitempty"`
	CreatedAt string              `json:"created_at" example:"2024-02-09T10:00:00Z"`
	UpdatedAt string              `json:"updated_at" example:"2024-02-09T10:00:00Z"`
}

func NewOrderProductJsonResponse(orderID uint64, productID uint64, quantity uint32) *OrderProductJsonResponse {
	orderProduct := &OrderProductJsonResponse{
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
	}

	return orderProduct
}

type OrderProductJsonPaginatedResponse struct {
	JsonPagination
	OrderProducts []OrderProductJsonResponse `json:"order_products"`
}
