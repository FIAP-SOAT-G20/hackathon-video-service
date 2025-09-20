package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/handler/response"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type HealthCheckHandler struct{}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) Register(router *gin.RouterGroup) {
	router.GET("", h.HealthCheck)
	router.GET("/", h.HealthCheck)
	router.GET("/readyz", h.HealthCheck)
	router.GET("/livez", h.HealthCheckLiveness)
}

// HealthCheck godoc
//
//	@Summary		Application Readiness
//	@Description	Checks application readiness
//	@Tags			health-check
//	@Produce		json
//	@Success		200	{object}	response.HealthCheckResponse
//	@Failure		500	{object}	string							"Internal server error"
//	@Failure		503	{object}	response.HealthCheckResponse	"Service Unavailable"
//	@Router			/health [GET]
//	@Router			/health/readyz [GET]
func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	cfg := config.LoadConfig()
	hc := &response.HealthCheckResponse{
		Status: response.HealthCheckStatusPass,
		Checks: map[string]response.HealthCheckVerifications{},
	}

	// Check databases based on configuration
	var hasFailures bool

	// Check PostgreSQL if using postgres engine
	if cfg.DBEngine == "postgres" || cfg.DBEngine == "postgresql" {
		postgresStatus := h.checkPostgreSQL(cfg)
		hc.Checks["postgres:status"] = postgresStatus
		if postgresStatus.Status == response.HealthCheckStatusFail {
			hasFailures = true
		}
	}

	// Check MongoDB if using documentdb/mongodb engine or as secondary database
	if cfg.DBEngine == "documentdb" || cfg.DBEngine == "mongodb" {
		mongoStatus := h.checkMongoDB(cfg)
		hc.Checks["mongodb:status"] = mongoStatus
		if mongoStatus.Status == response.HealthCheckStatusFail {
			hasFailures = true
		}
	}

	// Determine overall status
	if hasFailures {
		hc.Status = response.HealthCheckStatusFail
		c.JSON(http.StatusServiceUnavailable, hc)
		return
	}

	c.JSON(http.StatusOK, hc)
}

// checkPostgreSQL performs PostgreSQL health check
func (h *HealthCheckHandler) checkPostgreSQL(cfg *config.Config) response.HealthCheckVerifications {
	check := response.HealthCheckVerifications{
		ComponentId: "db:postgres",
		Status:      response.HealthCheckStatusPass,
		Time:        time.Now(),
	}

	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		check.Status = response.HealthCheckStatusFail
		return check
	}
	defer func() {
		_ = db.Close() // Ignore error for cleanup operation
	}()

	if err := db.Ping(); err != nil {
		check.Status = response.HealthCheckStatusFail
	}

	return check
}

// checkMongoDB performs MongoDB health check
func (h *HealthCheckHandler) checkMongoDB(cfg *config.Config) response.HealthCheckVerifications {
	check := response.HealthCheckVerifications{
		ComponentId: "db:mongodb",
		Status:      response.HealthCheckStatusPass,
		Time:        time.Now(),
	}

	// Set up MongoDB client options
	clientOptions := options.Client().ApplyURI(cfg.DocumentDBURI)
	clientOptions.SetConnectTimeout(10 * time.Second)
	clientOptions.SetServerSelectionTimeout(10 * time.Second)

	// Create MongoDB client
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		check.Status = response.HealthCheckStatusFail
		return check
	}
	defer func() {
		_ = client.Disconnect(context.Background()) // Ignore error for cleanup operation
	}()

	// Ping MongoDB to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		check.Status = response.HealthCheckStatusFail
	}

	return check
}

// HealthCheckLiveness godoc
//
//	@Summary		Application Liveness
//	@Description	Checks application liveness
//	@Tags			health-check
//	@Produce		json
//	@Success		200	{object}	response.HealthCheckLivenessResponse
//	@Router			/health/livez [GET]
func (h *HealthCheckHandler) HealthCheckLiveness(c *gin.Context) {
	hc := &response.HealthCheckLivenessResponse{
		Status: "ok",
	}

	c.JSON(http.StatusOK, hc)
}
