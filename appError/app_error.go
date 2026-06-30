package appError

import "fmt"

const (
	ErrNotFound       = "NOT_FOUND"
	ErrDuplicate      = "DUPLICATE"
	ErrFieldRequired  = "FIELD_REQUIRED"
	ErrCheckViolation = "CHECK_VIOLATION"
	ErrInvalidInput   = "INVALID_INPUT"
	ErrInternal       = "INTERNAL"
)

type AppError struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func New(status int, code string, message string, err error) *AppError {
	return &AppError{Status: status, Code: code, Message: message, Err: err}
}
