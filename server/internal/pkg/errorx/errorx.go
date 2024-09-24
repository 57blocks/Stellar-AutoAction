package errorx

import (
	"fmt"

	"github.com/57blocks/auto-action/server/internal/third-party/logx"
	
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

// BadRequest returns an error with status 400 and message.
func BadRequest(msg string) error {
	return fmt.Errorf("%w", newErr(http.StatusBadRequest, 400, msg))
}

// Unauthorized returns an error with status 400 and message.
func Unauthorized() error {
	return fmt.Errorf("%w", newErr(http.StatusUnauthorized, 401, "request unauthorized"))
}

// UnauthorizedWithMsg returns an error with status 400 and message.
func UnauthorizedWithMsg(msg string) error {
	return fmt.Errorf("%w", newErr(http.StatusUnauthorized, 401, msg))
}

// Forbidden returns an error with status 400 and message.
func Forbidden() error {
	return fmt.Errorf("%w", newErr(http.StatusForbidden, 403, "request forbidden"))
}

// ForbiddenWithMsg returns an error with status 400 and message.
func ForbiddenWithMsg(msg string) error {
	return fmt.Errorf("%w", newErr(http.StatusForbidden, 403, msg))
}

// NotFound returns an error with status 404 and message.
func NotFound(msg string) error {
	return fmt.Errorf("%w", newErr(http.StatusNotFound, 404, msg))
}

// Internal returns an error with status 404 and message.
func Internal(msg string) error {
	logx.Logger.ERROR(msg)

	return fmt.Errorf("%w", newErr(http.StatusInternalServerError, 500, msg))
}

func GinContextConv() error {
	return fmt.Errorf("%w", newErr(
		http.StatusInternalServerError,
		500,
		"convert context.Context to gin.Context failed",
	))
}

func AmazonConfig(msg string) error {
	logx.Logger.ERROR(msg)

	return fmt.Errorf("%w", newErr(
		http.StatusInternalServerError,
		500,
		fmt.Sprintf("failed to config Amazon: %s", msg),
	))
}
