package httpclient

import (
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/logger"
	"github.com/go-resty/resty/v2"
)

type HTTPClient struct {
	*resty.Client
}

func NewRestyClient(cfg *config.Config, logger *logger.Logger) *HTTPClient {
	httpCLient := resty.New().
		SetTimeout(10*1000*1000*1000). // 10 seconds
		SetRetryCount(2).
		// SetHeader("Authorization", fmt.Sprintf("Bearer %s", cfg.MercadoPagoToken)).
		SetHeader("Content-Type", "application/json")

	logger.Info("resty client created")

	return &HTTPClient{httpCLient}
}
