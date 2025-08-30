package domain

var (
	ErrConflict           = "data conflicts with existing data"
	ErrNotFound           = "data not found"
	ErrUnauthorized       = "unauthorized"
	ErrInvalidParam       = "invalid parameter"
	ErrInvalidQueryParams = "invalid query parameters"
	ErrInvalidBody        = "invalid body"

	ErrInvalidToken      = "access token is invalid"
	ErrMissingAuthHeader = "authorization header is required"
	ErrInvalidAuthHeader = "invalid authorization header format"

	ErrOrderInvalidStatusTransition = "invalid status transition"
	ErrOrderWithoutProducts         = "order without products"
	ErrProductIsMandatory           = "product is mandatory"
	ErrStaffIdIsMandatory           = "staff is mandatory"
	ErrOrderIsMandatory             = "order is mandatory"
	ErrOrderIsNotOpen               = "order is not on status open"
	ErrRoleInvalid                  = "invalid role"
	ErrStatusIsMandatory            = "status is mandatory"

	ErrPageMustBeGreaterThanZero = "page must be greater than zero"
	ErrLimitMustBeBetween1And100 = "limit must be between 1 and 100"

	ErrInternalError   = "internal server error"
	ErrUnknownError    = "unknown error"
	ErrValidationError = "validation error"
	ErrInvalidInput    = "invalid input"

	ErrFailedToCreatePaymentExternal = "failed to create payment external"
)

type ValidationError struct {
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type InternalError struct {
	Message string
	Err     error
}

func (e *InternalError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{
		Message: ErrValidationError,
		Err:     err,
	}
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}

func NewInternalError(err error) *InternalError {
	return &InternalError{
		Message: ErrInternalError,
		Err:     err,
	}
}

func NewInvalidInputError(message string) *InvalidInputError {
	return &InvalidInputError{
		Message: message,
	}
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{
		Message: message,
	}
}
