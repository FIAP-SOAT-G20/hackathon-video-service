package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/handler/response"

	"github.com/gin-gonic/gin"
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
		Checks: map[string]response.HealthCheckVerifications{
			"postgres:status": {
				ComponentId: "db:postgres",
				Status:      response.HealthCheckStatusPass,
				Time:        time.Now(),
			},
		},
	}

	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		_ = c.Error(err)
		return
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			_ = c.Error(closeErr)
		}
	}()

	if db.Ping() != nil {
		hc.Status = response.HealthCheckStatusFail
		hc.Checks["postgres:status"] = response.HealthCheckVerifications{
			ComponentId: "db:postgres",
			Status:      response.HealthCheckStatusFail,
			Time:        time.Now(),
		}
		c.JSON(http.StatusServiceUnavailable, hc)
		return
	}

	c.JSON(http.StatusOK, hc)
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
