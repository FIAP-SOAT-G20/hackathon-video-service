package request

type ListCategoriesQueryRequest struct {
	Name  string `form:"name" example:"Foods"`
	Page  int    `form:"page,default=1" example:"1"`
	Limit int    `form:"limit,default=10" example:"10"`
}

type CreateCategoryBodyRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100" example:"Foods"`
}

type GetCategoryUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type UpdateCategoryUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type UpdateCategoryBodyRequest struct {
	Name string `json:"name" binding:"omitempty,required" example:"Beverages"`
}

type DeleteCategoryUriRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}
