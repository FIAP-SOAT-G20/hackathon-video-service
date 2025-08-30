package handler

import (
	valueobject "github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/value_object"
	"github.com/go-playground/validator/v10"
)

func OrderStatusValidator(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	return valueobject.IsValidOrderStatus(status)
}
