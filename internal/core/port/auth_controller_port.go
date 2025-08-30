package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

type AuthController interface {
	Authenticate(ctx context.Context, presenter Presenter, input dto.AuthenticateInput) ([]byte, error)
}
