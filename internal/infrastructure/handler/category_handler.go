package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/adapter/presenter"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/handler/request"
)

type CategoryHandler struct {
	controller port.CategoryController
}

func NewCategoryHandler(controller port.CategoryController) *CategoryHandler {
	return &CategoryHandler{controller}
}

func (h *CategoryHandler) Register(router *gin.RouterGroup) {
	router.GET("", h.List)
	router.POST("", h.Create)
	router.GET("/:id", h.Get)
	router.PUT("/:id", h.Update)
	router.DELETE("/:id", h.Delete)
}

// List godoc
//
//	@Summary		List categories
//	@Description	List all categories
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int										false	"Page number"		default(1)
//	@Param			limit	query		int										false	"Items per page"	default(10)
//	@Success		200		{object}	presenter.CategoryJsonPaginatedResponse	"OK"
//	@Failure		400		{object}	middleware.ErrorJsonResponse			"Bad Request"
//	@Failure		500		{object}	middleware.ErrorJsonResponse			"Internal Server Error"
//	@Router			/categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	var query request.ListCategoriesQueryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidQueryParams))
		return
	}

	input := dto.ListCategoriesInput{
		Name:  query.Name,
		Page:  query.Page,
		Limit: query.Limit,
	}

	output, err := h.controller.List(
		c.Request.Context(),
		presenter.NewCategoryJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}

// Create godoc
//
//	@Summary		Create category
//	@Description	Creates a new category
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Param			category	body		request.CreateCategoryBodyRequest	true	"Category data"
//	@Success		201			{object}	presenter.CategoryJsonResponse		"Created"
//	@Failure		400			{object}	middleware.ErrorJsonResponse		"Bad Request"
//	@Failure		500			{object}	middleware.ErrorJsonResponse		"Internal Server Error"
//	@Router			/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var body request.CreateCategoryBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidBody))
		return
	}

	input := dto.CreateCategoryInput{
		Name: body.Name,
	}

	output, err := h.controller.Create(
		c.Request.Context(),
		presenter.NewCategoryJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusCreated, "application/json", output)
}

// Get godoc
//
//	@Summary		Get category
//	@Description	Search for a category by ID
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int								true	"Category ID"
//	@Success		200	{object}	presenter.CategoryJsonResponse	"OK"
//	@Failure		400	{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		404	{object}	middleware.ErrorJsonResponse	"Not Found"
//	@Failure		500	{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/categories/{id} [get]
func (h *CategoryHandler) Get(c *gin.Context) {
	var uri request.GetCategoryUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	input := dto.GetCategoryInput{
		ID: uri.ID,
	}

	output, err := h.controller.Get(
		c.Request.Context(),
		presenter.NewCategoryJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}

// Update godoc
//
//	@Summary		Update category
//	@Description	Update an existing category
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int									true	"Category ID"
//	@Param			category	body		request.UpdateCategoryBodyRequest	true	"Category data"
//	@Success		200			{object}	presenter.CategoryJsonResponse		"OK"
//	@Failure		400			{object}	middleware.ErrorJsonResponse		"Bad Request"
//	@Failure		404			{object}	middleware.ErrorJsonResponse		"Not Found"
//	@Failure		500			{object}	middleware.ErrorJsonResponse		"Internal Server Error"
//	@Router			/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	var uri request.UpdateCategoryUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	var body request.UpdateCategoryBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidBody))
		return
	}

	input := dto.UpdateCategoryInput{
		ID:   uri.ID,
		Name: body.Name,
	}

	output, err := h.controller.Update(
		c.Request.Context(),
		presenter.NewCategoryJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}

// Delete godoc
//
//	@Summary		Delete category
//	@Description	Deletes a category by ID
//	@Tags			category
//	@Produce		json
//	@Param			id	path		int								true	"Category ID"
//	@Success		200	{object}	presenter.CategoryJsonResponse	"OK"
//	@Failure		400	{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		404	{object}	middleware.ErrorJsonResponse	"Not Found"
//	@Failure		500	{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	var uri request.DeleteCategoryUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	input := dto.DeleteCategoryInput{
		ID: uri.ID,
	}

	output, err := h.controller.Delete(
		c.Request.Context(),
		presenter.NewCategoryJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}
