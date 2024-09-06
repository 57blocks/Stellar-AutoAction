package errorx

import (
	"fmt"
	"net/http"
)

type Errorx struct {
	StatusF int    `json:"status"`
	CodeF   int    `json:"code"`
	Message string `json:"message"`
}

func (e *Errorx) Status() int {
	return e.StatusF
}

func (e *Errorx) Code() int {
	return e.CodeF
}

func (e *Errorx) Error() string {
	return e.Message
}

func newErr(status, code int, message string) *Errorx {
	return &Errorx{
		StatusF: status,
		CodeF:   code,
		Message: message,
	}
}

// BadRequestErr returns an error with status 400 and message.
func BadRequestErr(msg string) error {
	return fmt.Errorf("%w", newErr(http.StatusBadRequest, 400, msg))
}

// UnauthorizedErr returns an error with status 400 and message.
func UnauthorizedErr() error {
	return fmt.Errorf("%w", newErr(http.StatusUnauthorized, 401, "request unauthorized"))
}

// UnauthorizedErrMsg returns an error with status 400 and message.
func UnauthorizedErrMsg(msg string) error {
	return fmt.Errorf("%w", newErr(http.StatusUnauthorized, 401, msg))
}

// ForbiddenErr returns an error with status 400 and message.
func ForbiddenErr() error {
	return fmt.Errorf("%w", newErr(http.StatusForbidden, 403, "request forbidden"))
}

// ForbiddenErrMsg returns an error with status 400 and message.
func ForbiddenErrMsg(msg string) error {
	return fmt.Errorf("%w", newErr(http.StatusForbidden, 403, msg))
}

// NotFoundErr returns an error with status 404 and message.
func NotFoundErr(msg string) error {
	return fmt.Errorf("%w", newErr(http.StatusNotFound, 404, msg))
}

// InternalErr returns an error with status 404 and message.
func InternalErr(msg string) error {
	return fmt.Errorf("%w", newErr(http.StatusInternalServerError, 500, msg))
}
