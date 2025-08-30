package request

type ListOrderProductsQueryRequest struct {
	OrderID   uint64 `form:"order_id,default=0" example:"1"`
	ProductID uint64 `form:"product_id,default=0" example:"1"`
	Page      int    `form:"page,default=1" example:"1"`
	Limit     int    `form:"limit,default=10" example:"10"`
}

type CreateOrderProductUriRequest struct {
	OrderID   uint64 `uri:"order_id" binding:"required"`
	ProductID uint64 `uri:"product_id" binding:"required"`
}

type CreateOrderProductBodyRequest struct {
	Quantity uint32 `json:"quantity" binding:"required" example:"1"`
}

type GetOrderProductUriRequest struct {
	OrderID   uint64 `uri:"order_id" binding:"required"`
	ProductID uint64 `uri:"product_id" binding:"required"`
}

type UpdateOrderProductUriRequest struct {
	OrderID   uint64 `uri:"order_id" binding:"required"`
	ProductID uint64 `uri:"product_id" binding:"required"`
}

type UpdateOrderProductBodyRequest struct {
	Quantity uint32 `json:"quantity" binding:"required" example:"1"`
}

type DeleteOrderProductUriRequest struct {
	OrderID   uint64 `uri:"order_id" binding:"required"`
	ProductID uint64 `uri:"product_id" binding:"required"`
}
