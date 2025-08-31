package handler

import (
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/go-playground/validator/v10"
)

func VideoStatusValidator(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	return valueobject.IsValidVideoStatus(status)
}
