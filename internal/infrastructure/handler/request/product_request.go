package request

type ListProductQueryRequest struct {
	Name       string `form:"name" example:"Product A"`
	CategoryID uint64 `form:"category_id" example:"1"`
	Page       int    `form:"page,default=1" example:"1"`
	Limit      int    `form:"limit,default=10" example:"10"`
}

type CreateProductBodyRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=100" example:"Product A"`
	Description string  `json:"description" binding:"max=500" example:"Product A description"`
	Price       float64 `json:"price" binding:"required,gt=0" example:"99.99"`
	CategoryID  uint64  `json:"category_id" binding:"required,gt=0" example:"1"`
}

// func (p *CreateProductRequest) Validate() error {
// 	return GetValidator().Struct(p)
// }

type GetProductUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type UpdateProductUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type UpdateProductBodyRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=100" example:"Product A"`
	Description string  `json:"description" binding:"max=500" example:"Product A description"`
	Price       float64 `json:"price" binding:"required,gt=0" example:"99.99"`
	CategoryID  uint64  `json:"category_id" binding:"required,gt=0" example:"1"`
}

type DeleteProductUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}
