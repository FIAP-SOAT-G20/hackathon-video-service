package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/adapter/presenter"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/handler/request"
)

type OrderHandler struct {
	controller port.OrderController
	jwtService port.JWTService
}

func NewOrderHandler(controller port.OrderController, jwtService port.JWTService) *OrderHandler {
	return &OrderHandler{controller: controller, jwtService: jwtService}
}

func (h *OrderHandler) Register(router *gin.RouterGroup) {
	//router.Use(middleware.JWTAuthMiddleware(h.jwtService))
	router.GET("", h.List)
	router.POST("", h.Create)
	router.GET("/:id", h.Get)
	router.PUT("/:id", h.Update)
	router.PATCH("/:id", h.UpdatePartial)
	router.DELETE("/:id", h.Delete)
}

// List godoc
//
//	@Summary		List orders (Reference TC-1 2.b.vi; TC-2 1.a.iv)
//	@Description	List all orders
//	@Description	## Order list is sorted by:
//	@Description	- **Status** in **descending** order (`READY` > `PREPARING` > `RECEIVED` > `PENDING` > `OPEN`)
//	@Description	- **Created date** (CreatedAt) in **ascending** order (oldest first)
//	@Description	Obs: Status CANCELLED and COMPLETED are not included in the list by default
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			customer_id		query		int										false	"Filter by customer ID"
//	@Param			status			query		string									false	"Filter by status (Accept many), options: <sub>OPEN, PENDING, RECEIVED, PREPARING, READY</sub>, ex: <sub>PENDING</sub> or <sub>OPEN,PENDING</sub>"
//	@Param			status_exclude	query		string									false	"Exclude by status (Accept many), options: <sub>NONE, OPEN, PENDING, RECEIVED, PREPARING, READY, CANCELLED, COMPLETED</sub>, ex: <sub>CANCELLED</sub> or <sub>CANCELLED,COMPLETED</sub> (default)"	default(CANCELLED,COMPLETED)
//	@Param			sort			query		string									false	"Sort by field (Accept many). Use `<field_name>:d` for descending, and the default order is ascending"																								default(status:d,created_at)
//	@Param			page			query		int										false	"Page number"																																														default(1)
//	@Param			limit			query		int										false	"Items per page"																																													default(10)
//	@Success		200				{object}	presenter.OrderJsonPaginatedResponse	"OK"
//	@Failure		400				{object}	middleware.ErrorJsonResponse			"Bad Request"
//	@Failure		500				{object}	middleware.ErrorJsonResponse			"Internal Server Error"
//	@Router			/orders [get]
func (h *OrderHandler) List(c *gin.Context) {
	var query request.ListOrdersQueryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	// Default sort
	if query.Sort == "" {
		query.Sort = "status:d,created_at"
	}

	// Default status_exclude
	var statusExclude []valueobject.OrderStatus
	if query.StatusExclude == "" {
		query.StatusExclude = "CANCELLED,COMPLETED"
	}

	// Convert status_exclude
	if strings.ToUpper(query.StatusExclude) != "NONE" {
		for _, s := range strings.Split(query.StatusExclude, ",") {
			statusExclude = append(statusExclude, valueobject.OrderStatus(strings.TrimSpace(s)))
		}
	}

	// Convert status
	var status []valueobject.OrderStatus
	if query.Status != "" {
		for _, s := range strings.Split(query.Status, ",") {
			orderStatus, ok := valueobject.ToOrderStatus(strings.TrimSpace(s))
			if !ok {
				_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
				return
			}
			status = append(status, orderStatus)
		}
	}

	input := dto.ListOrdersInput{
		CustomerID:    query.CustomerID,
		Status:        status,
		StatusExclude: statusExclude,
		Page:          query.Page,
		Limit:         query.Limit,
		Sort:          query.Sort,
	}

	output, err := h.controller.List(
		c.Request.Context(),
		presenter.NewOrderJsonPresenter(),
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
//	@Summary		Create order
//	@Description	Creates a new order
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			order	body		request.CreateOrderBodyRequest	true	"Order data"
//	@Success		201		{object}	presenter.OrderJsonResponse		"Created"
//	@Failure		400		{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		500		{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var body request.CreateOrderBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidBody))
		return
	}

	input := dto.CreateOrderInput{
		CustomerID: body.CustomerID,
	}

	output, err := h.controller.Create(
		c.Request.Context(),
		presenter.NewOrderJsonPresenter(),
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
//	@Summary		Get order
//	@Description	Search for a order by ID
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int								true	"Order ID"
//	@Success		200	{object}	presenter.OrderJsonResponse		"OK"
//	@Failure		400	{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		404	{object}	middleware.ErrorJsonResponse	"Not Found"
//	@Failure		500	{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/orders/{id} [get]
func (h *OrderHandler) Get(c *gin.Context) {
	var uri request.GetOrderUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	input := dto.GetOrderInput{
		ID: uri.ID,
	}

	output, err := h.controller.Get(
		c.Request.Context(),
		presenter.NewOrderJsonPresenter(),
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
//	@Summary		Update order
//	@Description	Update an existing order
//	@Description	The status are: **OPEN**, **CANCELLED**, **PENDING**, **RECEIVED**, **PREPARING**, **READY**, **COMPLETED**
//	@Description	## Transition of status:
//	@Description	- OPEN      -> CANCELLED || PENDING
//	@Description	- CANCELLED -> {},
//	@Description	- PENDING   -> OPEN || RECEIVED
//	@Description	- RECEIVED  -> PREPARING
//	@Description	- PREPARING -> READY
//	@Description	- READY     -> COMPLETED
//	@Description	- COMPLETED -> {}
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Order ID"
//	@Param			order	body		request.UpdateOrderBodyRequest	true	"Order data"
//	@Success		200		{object}	presenter.OrderJsonResponse		"OK"
//	@Failure		400		{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		404		{object}	middleware.ErrorJsonResponse	"Not Found"
//	@Failure		500		{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/orders/{id} [put]
func (h *OrderHandler) Update(c *gin.Context) {
	var uri request.UpdateOrderUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	var body request.UpdateOrderBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidBody))
		return
	}

	input := dto.UpdateOrderInput{
		ID:         uri.ID,
		CustomerID: body.CustomerID,
		Status:     body.Status,
		StaffID:    body.StaffID,
	}

	output, err := h.controller.Update(
		c.Request.Context(),
		presenter.NewOrderJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}

// UpdatePartial godoc
//
//	@Summary		Partial update order (Reference TC-2 1.a.v)
//	@Description	Partially updates an existing order
//	@Description	The status are: **OPEN**, **CANCELLED**, **PENDING**, **RECEIVED**, **PREPARING**, **READY**, **COMPLETED**
//	@Description	## Transition of status:
//	@Description	- OPEN      -> CANCELLED || PENDING
//	@Description	- CANCELLED -> {},
//	@Description	- PENDING   -> OPEN || RECEIVED
//	@Description	- RECEIVED  -> PREPARING
//	@Description	- PREPARING -> READY
//	@Description	- READY     -> COMPLETED
//	@Description	- COMPLETED -> {}
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"Order ID"
//	@Param			order	body		request.UpdateOrderPartilRequest	true	"Order data"
//	@Success		200		{object}	presenter.OrderJsonResponse			"OK"
//	@Failure		400		{object}	middleware.ErrorJsonResponse		"Bad Request"
//	@Failure		404		{object}	middleware.ErrorJsonResponse		"Not Found"
//	@Failure		500		{object}	middleware.ErrorJsonResponse		"Internal Server Error"
//	@Router			/orders/{id} [patch]
func (h *OrderHandler) UpdatePartial(c *gin.Context) {
	var uri request.UpdateOrderUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	fmt.Println("URI ID:", uri.ID)

	var body request.UpdateOrderPartilBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println(err)
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidBody))
		return
	}

	fmt.Println("Body:", body)

	input := dto.UpdateOrderInput{
		ID:         uri.ID,
		CustomerID: body.CustomerID,
		Status:     body.Status,
		StaffID:    body.StaffID,
	}

	output, err := h.controller.Update(
		c.Request.Context(),
		presenter.NewOrderJsonPresenter(),
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
//	@Summary		Delete order
//	@Description	Deletes a order by ID
//	@Tags			orders
//	@Produce		json
//	@Param			id	path		int								true	"Order ID"
//	@Success		200	{object}	presenter.OrderJsonResponse		"OK"
//	@Failure		400	{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		404	{object}	middleware.ErrorJsonResponse	"Not Found"
//	@Failure		500	{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/orders/{id} [delete]
func (h *OrderHandler) Delete(c *gin.Context) {
	var uri request.DeleteOrderUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	input := dto.DeleteOrderInput{
		ID: uri.ID,
	}

	output, err := h.controller.Delete(
		c.Request.Context(),
		presenter.NewOrderJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}
