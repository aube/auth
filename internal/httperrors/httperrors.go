package httperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code    int
	Err     error
	Message string
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("%d - %s", he.Code, he.Message)
}

func NewHTTPError(code int, message string) error {
	return &HTTPError{
		Code:    code,
		Message: message,
		Err:     errors.New(message),
	}
}

func NewServerError(err error) error {
	return NewHTTPError(500, fmt.Sprint(err))
}

func NewRecordNotFound() error {
	return NewHTTPError(404, "Record not found")
}

func NewValidationError(validationError error) error {
	return NewHTTPError(http.StatusBadRequest, fmt.Sprint(validationError))
}

func NewConflictError() error {
	return NewHTTPError(http.StatusConflict, "Conflict") // 409
}

func NewLoginFailed() error {
	return NewHTTPError(http.StatusForbidden, "Login failed") // 403
}

func NewAccessDenied() error {
	return NewHTTPError(http.StatusForbidden, "Access denied") // 403
}

func NewAlreadyUploadedByMeError() error {
	return NewHTTPError(http.StatusOK, "Already uploaded") // 200
}

func NewAlreadyUploadedAnotherError() error {
	return NewHTTPError(http.StatusConflict, "Already uploaded") // 200
}

func NewNotEnoughMoneyError() error {
	return NewHTTPError(402, "I need more gold!")
}

func NewOrderNumberError() error {
	return NewHTTPError(http.StatusUnprocessableEntity, "Order number incorrect") // 422
}
