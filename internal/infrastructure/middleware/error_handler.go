package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/logger"
)

type ErrorJsonResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}

type ErrorXmlResponse struct {
	Code    int    `xml:"code" example:"400"`
	Message string `xml:"message" example:"Bad Request"`
}

func ErrorHandler(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Execute all the handlers

		// If there are errors, handle the last one
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err, logger)
		}
	}
}

func handleError(c *gin.Context, err error, logger *logger.Logger) {
	switch e := err.(type) {
	case *domain.ValidationError:
		setResponse(c, http.StatusBadRequest, e.Error())
		logWarning(logger, domain.ErrValidationError, e, c.Request)

	case *domain.NotFoundError:
		setResponse(c, http.StatusNotFound, e.Error())
		logWarning(logger, domain.ErrNotFound, e, c.Request)

	case *domain.InvalidInputError:
		setResponse(c, http.StatusBadRequest, e.Error())
		logWarning(logger, domain.ErrInvalidInput, e, c.Request)

	case *domain.UnauthorizedError:
		setResponse(c, http.StatusUnauthorized, e.Error())
		logWarning(logger, domain.ErrUnauthorized, e, c.Request)

	case *domain.InternalError:
		setResponse(c, http.StatusInternalServerError, domain.ErrInternalError)
		logError(logger, domain.ErrInternalError, e, c.Request)

	default:
		setResponse(c, http.StatusInternalServerError, domain.ErrInternalError)
		logError(logger, domain.ErrUnknownError, err, c.Request)
	}
}

func setResponse(c *gin.Context, status int, message string) {
	if c.GetHeader("Accept") == "text/xml" {
		c.XML(status, ErrorXmlResponse{
			Code:    status,
			Message: message,
		})
		return
	}

	c.JSON(status, ErrorJsonResponse{
		Code:    status,
		Message: message,
	})
}

func logError(logger *logger.Logger, msg string, err error, req *http.Request) {
	logger.Error(msg,
		"error", err.Error(),
		"path", req.URL.Path,
		"method", req.Method,
	)
}

func logWarning(logger *logger.Logger, msg string, err error, req *http.Request) {
	logger.Warn(msg,
		"error", err.Error(),
		"path", req.URL.Path,
		"method", req.Method,
	)
}
