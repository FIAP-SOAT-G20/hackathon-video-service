package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
)

// AuthUseCase defines the authentication use case interface
type AuthUseCase interface {
	// Authenticate authenticates a customer by CPF and return the token
	Authenticate(ctx context.Context, input dto.AuthenticateInput) (string, error)
}
