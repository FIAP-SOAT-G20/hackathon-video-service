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

type ProductHandler struct {
	controller port.ProductController
}

func NewProductHandler(controller port.ProductController) *ProductHandler {
	return &ProductHandler{controller: controller}
}

func (h *ProductHandler) Register(router *gin.RouterGroup) {
	router.GET("", h.List)
	router.POST("", h.Create)
	router.GET("/:id", h.Get)
	router.PUT("/:id", h.Update)
	router.DELETE("/:id", h.Delete)
}

// List godoc
//
//	@Summary		List products (Reference TC-1 2.b.iv)
//	@Description	List all products
//	@Description	Response can return JSON or XML format (Accept header: application/json or text/xml)
//	@Tags			products
//	@Accept			json
//	@Produce		json,xml
//	@Param			name		query		string									false	"Filter by name"
//	@Param			category_id	query		int										false	"Filter by category ID"
//	@Param			page		query		int										false	"Page number"		default(1)
//	@Param			limit		query		int										false	"Items per page"	default(10)
//	@Success		200			{object}	presenter.ProductJsonPaginatedResponse	"OK"
//	@Failure		400			{object}	middleware.ErrorJsonResponse			"Bad Request"
//	@Failure		500			{object}	middleware.ErrorJsonResponse			"Internal Server Error"
//	@Router			/products [get]
func (h *ProductHandler) List(c *gin.Context) {
	var query request.ListProductQueryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidQueryParams))
		return
	}

	input := dto.ListProductsInput{
		Name:       query.Name,
		CategoryID: query.CategoryID,
		Page:       query.Page,
		Limit:      query.Limit,
	}

	p, contentType := selectOutputConfigs(c.GetHeader("Accept"))

	output, err := h.controller.List(c.Request.Context(), p, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, contentType, output)
}

// Create godoc
//
//	@Summary		Create product (Reference TC-1 2.b.iii)
//	@Description	Creates a new product
//	@Description	Response can return JSON or XML format (Accept header: application/json or text/xml)
//	@Tags			products
//	@Accept			json
//	@Produce		json,xml
//	@Param			product	body		request.CreateProductBodyRequest	true	"Product data"
//	@Success		201		{object}	presenter.ProductJsonResponse		"Created"
//	@Failure		400		{object}	middleware.ErrorJsonResponse		"Bad Request"
//	@Failure		500		{object}	middleware.ErrorJsonResponse		"Internal Server Error"
//	@Router			/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var body request.CreateProductBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidBody))
		return
	}

	input := dto.CreateProductInput{
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		CategoryID:  body.CategoryID,
	}

	p, contentType := selectOutputConfigs(c.GetHeader("Accept"))

	output, err := h.controller.Create(c.Request.Context(), p, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusCreated, contentType, output)
}

// Get godoc
//
//	@Summary		Get product
//	@Description	Search for a product by ID
//	@Description	Response can return JSON or XML format (Accept header: application/json or text/xml)
//	@Tags			products
//	@Accept			json
//	@Produce		json,xml
//	@Param			id	path		int								true	"Product ID"
//	@Success		200	{object}	presenter.ProductJsonResponse	"OK"
//	@Failure		400	{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		404	{object}	middleware.ErrorJsonResponse	"Not Found"
//	@Failure		500	{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/products/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	var uri request.GetProductUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	input := dto.GetProductInput{
		ID: uri.ID,
	}

	p, contentType := selectOutputConfigs(c.GetHeader("Accept"))

	output, err := h.controller.Get(c.Request.Context(), p, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, contentType, output)
}

// Update godoc
//
//	@Summary		Update product (Reference TC-1 2.b.iii)
//	@Description	Update an existing product
//	@Description	Response can return JSON or XML format (Accept header: application/json or text/xml)
//	@Tags			products
//	@Accept			json
//	@Produce		json,xml
//	@Param			id		path		int									true	"Product ID"
//	@Param			product	body		request.UpdateProductBodyRequest	true	"Product data"
//	@Success		200		{object}	presenter.ProductJsonResponse		"OK"
//	@Failure		400		{object}	middleware.ErrorJsonResponse		"Bad Request"
//	@Failure		404		{object}	middleware.ErrorJsonResponse		"Not Found"
//	@Failure		500		{object}	middleware.ErrorJsonResponse		"Internal Server Error"
//	@Router			/products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	var uri request.UpdateProductUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	var body request.UpdateProductBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidBody))
		return
	}

	input := dto.UpdateProductInput{
		ID:          uri.ID,
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		CategoryID:  body.CategoryID,
	}

	p, contentType := selectOutputConfigs(c.GetHeader("Accept"))

	output, err := h.controller.Update(c.Request.Context(), p, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, contentType, output)
}

// Delete godoc
//
//	@Summary		Delete product (Reference TC-1 2.b.iii)
//	@Description	Deletes a product by ID
//	@Description	Response can return JSON or XML format (Accept header: application/json or text/xml)
//	@Tags			products
//	@Accept			json
//	@Produce		json,xml
//	@Param			id	path		int								true	"Product ID"
//	@Success		200	{object}	presenter.ProductJsonResponse	"OK"
//	@Failure		400	{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		404	{object}	middleware.ErrorJsonResponse	"Not Found"
//	@Failure		500	{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	var uri request.DeleteProductUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	input := dto.DeleteProductInput{
		ID: uri.ID,
	}

	p, contentType := selectOutputConfigs(c.GetHeader("Accept"))

	output, err := h.controller.Delete(c.Request.Context(), p, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, contentType, output)
}

func selectOutputConfigs(acceptHeader string) (port.Presenter, string) {
	if acceptHeader == "text/xml" {
		return presenter.NewProductXmlPresenter(), acceptHeader
	}
	return presenter.NewProductJsonPresenter(), "application/json"
}
