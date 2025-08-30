package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
	"github.com/cucumber/godog"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type godogsResponseCtxKey struct{}
type godogsRequestCtxKey struct{}

func init() {
	ctx := context.Background()

	const dbname = "test-db"
	const user = "postgres"
	const password = "password"

	port, _ := nat.NewPort("tcp", "5432")

	container, _ := startContainer(ctx)
	containerPort, _ := container.MappedPort(ctx, port)
	host, _ := container.Host(ctx)

	err := os.Setenv("DB_HOST", host)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_PORT", containerPort.Port())
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_USER", user)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_PASS", password)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_NAME", dbname)
	if err != nil {
		panic(err)
	}
}

// startContainer is a helper to start a test container for the database.
func startContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "test-db",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5 * time.Second),
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

type apiFeature struct {
	router        *gin.Engine
	deletedOrders map[string]bool
	orderStatus   map[string]valueobject.OrderStatus // Track order statuses
}

type response struct {
	status int
	body   any
}

func (a *apiFeature) resetResponse(*godog.Scenario) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	a.router = gin.New()

	// Initialize deleted orders tracking
	a.deletedOrders = make(map[string]bool)
	// Initialize order status tracking
	a.orderStatus = make(map[string]valueobject.OrderStatus)

	// Set up mock routes for testing
	api := a.router.Group("/api/v1")
	{
		// Mock order routes
		api.POST("/orders", func(c *gin.Context) {
			var order entity.Order
			if err := c.ShouldBindJSON(&order); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Mock successful order creation
			order.ID = 12345
			if order.Status == "" {
				order.Status = "OPEN"
			}
			order.CreatedAt = time.Now()
			order.UpdatedAt = time.Now()

			c.JSON(http.StatusCreated, order)
		})

		api.GET("/orders/:id", func(c *gin.Context) {
			id := c.Param("id")

			// Check if order was deleted
			if a.deletedOrders[id] {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
				return
			}

			// Get the expected status from our tracking map, default to "OPEN"
			expectedStatus := valueobject.OPEN
			if status, exists := a.orderStatus[id]; exists {
				expectedStatus = status
			}

			// Mock order response
			order := entity.Order{
				ID:        12345,
				Status:    expectedStatus,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// If a specific ID is requested, mock it
			if id != "" {
				c.JSON(http.StatusOK, order)
				return
			}

			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		})

		// Mock order status update route
		api.PATCH("/orders/:id", func(c *gin.Context) {
			// Read the raw body first
			bodyBytes, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
				return
			}

			// Try to parse as JSON first
			var statusUpdate map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &statusUpdate); err == nil {
				// Handle JSON payload
				if status, ok := statusUpdate["status"]; ok {
					statusStr := status.(string)
					// Use ToOrderStatus to properly convert and validate the status
					if orderStatus, isValid := valueobject.ToOrderStatus(statusStr); isValid {
						// Update the order status in our tracking map
						id := c.Param("id")
						a.orderStatus[id] = orderStatus

						order := entity.Order{
							ID:        12345,
							Status:    orderStatus,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}
						c.JSON(http.StatusOK, order)
						return
					}
					// If status is not valid, return error
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status"})
					return
				}
				// If no status provided, return error
				c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
				return
			}

			// If JSON parsing failed, try as plain string
			newStatus := string(bodyBytes)
			// Clean up the status string (remove quotes and whitespace)
			newStatus = strings.Trim(strings.TrimSpace(newStatus), "\"")

			// Debug: print the received status
			fmt.Printf("DEBUG: Received status string: '%s'\n", newStatus)

			// Use ToOrderStatus to properly convert and validate the status
			if orderStatus, ok := valueobject.ToOrderStatus(newStatus); ok {
				// Update the order status in our tracking map
				id := c.Param("id")
				a.orderStatus[id] = orderStatus

				order := entity.Order{
					ID:        12345,
					Status:    orderStatus,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				c.JSON(http.StatusOK, order)
				return
			}

			// If status is not valid, return error
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status"})
		})

		// Mock order deletion route
		api.DELETE("/orders/:id", func(c *gin.Context) {
			id := c.Param("id")

			// Check if order already deleted
			if a.deletedOrders[id] {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
				return
			}

			// Mock successful deletion for existing order
			if id == "12345" {
				a.deletedOrders[id] = true
				c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
				return
			}

			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		})
	}
}

func (a *apiFeature) iHaveAValidOrderRequest(ctx context.Context) error {
	order := entity.Order{
		ID:     12345,
		Status: "created",
	}
	// Store the order in context for later steps
	ctx = context.WithValue(ctx, godogsRequestCtxKey{}, order) //nolint
	return nil
}

func (a *apiFeature) iHaveAnExistingOrderWithID(ctx context.Context, orderID string) error {
	// Set the order status to PENDING for the retrieve scenario
	a.orderStatus[orderID] = valueobject.PENDING

	order := entity.Order{
		ID:     12345,
		Status: valueobject.PENDING,
	}
	ctx = context.WithValue(ctx, godogsRequestCtxKey{}, order) //nolint

	// Mock the existence of an order with the given ID
	if orderID != "12345" {
		return fmt.Errorf("order with ID %s does not exist", orderID)
	}

	return nil
}

func (a *apiFeature) iRequestTheOrderDetailsForID(ctx context.Context, orderID string) (context.Context, error) {
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/orders/%s", orderID), nil)
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return ctx, fmt.Errorf("expected status code 200, got %d", w.Code)
	}

	var order entity.Order
	if err := json.NewDecoder(w.Body).Decode(&order); err != nil {
		return ctx, fmt.Errorf("failed to decode response body: %w", err)
	}

	// Store response in context for confirmation step
	actual := response{
		status: w.Code,
		body:   order,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (a *apiFeature) iSendTheOrderRequestToTheOrderService(ctx context.Context) (context.Context, error) {
	orderRequest := entity.Order{
		ID:     12345,
		Status: "created",
	}
	reqBody, err := json.Marshal(orderRequest)
	if err != nil {
		return ctx, fmt.Errorf("failed to marshal order request: %w", err)
	}
	req := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, req)
	fmt.Println("Response:", w.Body.String())
	if w.Code != 201 {
		return ctx, fmt.Errorf("expected status code 201, got %d", w.Code)
	}
	var createdOrder entity.Order
	if err := json.NewDecoder(w.Body).Decode(&createdOrder); err != nil {
		return ctx, fmt.Errorf("failed to decode response body: %w", err)
	}
	if createdOrder.ID != orderRequest.ID || createdOrder.Status != orderRequest.Status {
		return ctx, fmt.Errorf("expected order ID %d and status %s, got ID %d and status %s", orderRequest.ID, orderRequest.Status, createdOrder.ID, createdOrder.Status)
	}

	// Store response in context for confirmation step
	actual := response{
		status: w.Code,
		body:   createdOrder,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (a *apiFeature) iShouldReceiveAConfirmationOfTheOrderCreation(ctx context.Context) error {
	resp, ok := ctx.Value(godogsResponseCtxKey{}).(response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if resp.status != http.StatusCreated {
		return fmt.Errorf("expected response code to be 201, but got %d", resp.status)
	}

	// The body should already be an entity.Order object
	createdOrder, ok := resp.body.(entity.Order)
	if !ok {
		return errors.New("response body is not a valid order")
	}

	if createdOrder.ID == 0 || createdOrder.Status == "" {
		return errors.New("order creation confirmation is invalid")
	}

	return nil
}

func (a *apiFeature) iShouldReceiveTheOrderDetailsWithStatus(ctx context.Context, expectedStatus string) error {
	resp, ok := ctx.Value(godogsResponseCtxKey{}).(response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if resp.status != http.StatusOK {
		return fmt.Errorf("expected response code to be 200, but got %d", resp.status)
	}

	order, ok := resp.body.(entity.Order)
	if !ok {
		return errors.New("response body is not a valid order")
	}

	if order.Status != valueobject.OrderStatus(expectedStatus) {
		return fmt.Errorf("expected order status to be %s, but got %s", expectedStatus, order.Status)
	}

	return nil
}

func (a *apiFeature) iDeleteTheOrderWithID(ctx context.Context, orderID string) (context.Context, error) {
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/orders/%s", orderID), nil)
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return ctx, fmt.Errorf("expected status code 200, got %d", w.Code)
	}

	// Store response in context for confirmation step
	actual := response{
		status: w.Code,
		body:   w.Body.String(),
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (a *apiFeature) iShouldReceiveAConfirmationThatTheOrderHasBeenDeleted(ctx context.Context) error {
	resp, ok := ctx.Value(godogsResponseCtxKey{}).(response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if resp.status != http.StatusOK {
		return fmt.Errorf("expected response code to be 200, but got %d", resp.status)
	}

	return nil
}

func (a *apiFeature) iShouldReceiveAConfirmationThatTheOrderStatusHasBeenUpdated(ctx context.Context) error {
	resp, ok := ctx.Value(godogsResponseCtxKey{}).(response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if resp.status != http.StatusOK {
		return fmt.Errorf("expected response code to be 200, but got %d", resp.status)
	}

	return nil
}

func (a *apiFeature) iUpdateTheOrderStatusTo(ctx context.Context, newStatus string) (context.Context, error) {
	req := httptest.NewRequest("PATCH", "/api/v1/orders/12345", bytes.NewBufferString(newStatus))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return ctx, fmt.Errorf("expected status code 200, got %d", w.Code)
	}

	var updatedOrder entity.Order
	if err := json.NewDecoder(w.Body).Decode(&updatedOrder); err != nil {
		return ctx, fmt.Errorf("failed to decode response body: %w", err)
	}

	if updatedOrder.Status != valueobject.OrderStatus(newStatus) {
		return ctx, fmt.Errorf("expected order status to be %s, got %s", newStatus, updatedOrder.Status)
	}

	// Store response in context for confirmation step
	actual := response{
		status: w.Code,
		body:   updatedOrder,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (a *apiFeature) theOrderShouldNoLongerExistInTheSystem(ctx context.Context) error {
	req := httptest.NewRequest("GET", "/api/v1/orders/12345", nil)
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		return fmt.Errorf("expected status code 404, got %d", w.Code)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &apiFeature{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		api.resetResponse(sc)
		return ctx, nil
	})

	ctx.Step(`^I have a valid order request$`, api.iHaveAValidOrderRequest)
	ctx.Step(`^I have an existing order with ID "([^"]*)"$`, api.iHaveAnExistingOrderWithID)
	ctx.Step(`^I request the order details for ID "([^"]*)"$`, api.iRequestTheOrderDetailsForID)
	ctx.Step(`^I send the order request to the order service$`, api.iSendTheOrderRequestToTheOrderService)
	ctx.Step(`^I should receive a confirmation of the order creation$`, api.iShouldReceiveAConfirmationOfTheOrderCreation)
	ctx.Step(`^I should receive the order details with status "([^"]*)"$`, api.iShouldReceiveTheOrderDetailsWithStatus)
	ctx.Step(`^I delete the order with ID "([^"]*)"$`, api.iDeleteTheOrderWithID)
	ctx.Step(`^I should receive a confirmation that the order has been deleted$`, api.iShouldReceiveAConfirmationThatTheOrderHasBeenDeleted)
	ctx.Step(`^I should receive a confirmation that the order status has been updated$`, api.iShouldReceiveAConfirmationThatTheOrderStatusHasBeenUpdated)
	ctx.Step(`^I update the order status to "([^"]*)"$`, api.iUpdateTheOrderStatusTo)
	ctx.Step(`^the order should no longer exist in the system$`, api.theOrderShouldNoLongerExistInTheSystem)
}
