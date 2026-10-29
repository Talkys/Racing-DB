package apperror

import (
	"fmt"
	"net/http"
)

// Error is an error caused by the client's request. Its message is safe to send back;
// any other error is treated as internal and answered with a generic 500.
type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func BadRequest(format string, v ...any) *Error {
	return &Error{Status: http.StatusBadRequest, Message: fmt.Sprintf(format, v...)}
}

func NotFound(format string, v ...any) *Error {
	return &Error{Status: http.StatusNotFound, Message: fmt.Sprintf(format, v...)}
}
