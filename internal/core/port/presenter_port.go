package port

import "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"

type Presenter interface {
	Present(dto.PresenterInput) ([]byte, error)
}
