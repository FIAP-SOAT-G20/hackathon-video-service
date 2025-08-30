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

type OrderProductHandler struct {
	controller port.OrderProductController
}

func NewOrderProductHandler(controller port.OrderProductController) *OrderProductHandler {
	return &OrderProductHandler{controller: controller}
}

func (h *OrderProductHandler) Register(router *gin.RouterGroup) {
	router.GET("", h.List)
	router.POST("/:order_id/:product_id", h.Create)
	router.GET("/:order_id/:product_id", h.Get)
	router.PUT("/:order_id/:product_id", h.Update)
	router.DELETE("/:order_id/:product_id", h.Delete)
}

// List godoc
//
//	@Summary		List order products
//	@Description	List all order products
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			order_id	query		string										false	"Filter by order ID"
//	@Param			page		query		int											false	"Page number"		default(1)
//	@Param			limit		query		int											false	"Items per page"	default(10)
//	@Success		200			{object}	presenter.OrderProductJsonPaginatedResponse	"OK"
//	@Failure		400			{object}	middleware.ErrorJsonResponse				"Bad Request"
//	@Failure		500			{object}	middleware.ErrorJsonResponse				"Internal Server Error"
//	@Router			/orders/products [get]
func (h *OrderProductHandler) List(c *gin.Context) {
	var query request.ListOrderProductsQueryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	input := dto.ListOrderProductsInput{
		OrderID:   query.OrderID,
		ProductID: query.ProductID,
		Page:      query.Page,
		Limit:     query.Limit,
	}

	output, err := h.controller.List(
		c.Request.Context(),
		presenter.NewOrderProductJsonPresenter(),
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
//	@Summary		Create an order product
//	@Description	Create an order product
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			order_id	path		int										true	"Order ID"
//	@Param			product_id	path		int										true	"Product ID"
//	@Param			order		body		request.CreateOrderProductBodyRequest	true	"OrderProduct data"
//	@Success		201			{object}	presenter.OrderProductJsonResponse		"Created"
//	@Failure		400			{object}	middleware.ErrorJsonResponse			"Bad Request"
//	@Router			/orders/products/{order_id}/{product_id} [post]
func (h *OrderProductHandler) Create(c *gin.Context) {
	var uri request.CreateOrderProductUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	var body request.CreateOrderProductBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidBody))
		return
	}

	input := dto.CreateOrderProductInput{
		OrderID:   uri.OrderID,
		ProductID: uri.ProductID,
		Quantity:  body.Quantity,
	}

	output, err := h.controller.Create(
		c.Request.Context(),
		presenter.NewOrderProductJsonPresenter(),
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
//	@Summary		Get an order product
//	@Description	Get an order product
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			order_id	path		int									true	"Order ID"
//	@Param			product_id	path		int									true	"Product ID"
//	@Success		200			{object}	presenter.OrderProductJsonResponse	"OK"
//	@Failure		400			{object}	middleware.ErrorJsonResponse		"Bad Request"
//	@Failure		404			{object}	middleware.ErrorJsonResponse		"Not Found"
//	@Failure		500			{object}	middleware.ErrorJsonResponse		"Internal Server Error"
//	@Router			/orders/products/{order_id}/{product_id} [get]
func (h *OrderProductHandler) Get(c *gin.Context) {
	var uri request.GetOrderProductUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	input := dto.GetOrderProductInput{
		OrderID:   uri.OrderID,
		ProductID: uri.ProductID,
	}

	output, err := h.controller.Get(
		c.Request.Context(),
		presenter.NewOrderProductJsonPresenter(),
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
//	@Summary		Update order product
//	@Description	Update an existing order product
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			order_id	path		int										true	"Order ID"
//	@Param			product_id	path		int										true	"Product ID"
//	@Param			order		body		request.UpdateOrderProductBodyRequest	true	"OrderProduct data"
//	@Success		200			{object}	presenter.OrderProductJsonResponse		"OK"
//	@Failure		400			{object}	middleware.ErrorJsonResponse			"Bad Request"
//	@Failure		404			{object}	middleware.ErrorJsonResponse			"Not Found"
//	@Failure		500			{object}	middleware.ErrorJsonResponse			"Internal Server Error"
//	@Router			/orders/products/{order_id}/{product_id} [put]
func (h *OrderProductHandler) Update(c *gin.Context) {
	var uri request.UpdateOrderProductUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	var body request.UpdateOrderProductBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidBody))
		return
	}

	input := dto.UpdateOrderProductInput{
		OrderID:   uri.OrderID,
		ProductID: uri.ProductID,
		Quantity:  body.Quantity,
	}

	output, err := h.controller.Update(
		c.Request.Context(),
		presenter.NewOrderProductJsonPresenter(),
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
//	@Summary		Delete order product
//	@Description	Deletes a order product by Order ID and Product ID
//	@Tags			orders
//	@Produce		json
//	@Param			order_id	path		int									true	"Order ID"
//	@Param			product_id	path		int									true	"Product ID"
//	@Success		200			{object}	presenter.OrderProductJsonResponse	"OK"
//	@Failure		400			{object}	middleware.ErrorJsonResponse		"Bad Request"
//	@Failure		404			{object}	middleware.ErrorJsonResponse		"Not Found"
//	@Failure		500			{object}	middleware.ErrorJsonResponse		"Internal Server Error"
//	@Router			/orders/products/{order_id}/{product_id} [delete]
func (h *OrderProductHandler) Delete(c *gin.Context) {
	var uri request.DeleteOrderProductUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	input := dto.DeleteOrderProductInput{
		OrderID:   uri.OrderID,
		ProductID: uri.ProductID,
	}

	output, err := h.controller.Delete(
		c.Request.Context(),
		presenter.NewOrderProductJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}
