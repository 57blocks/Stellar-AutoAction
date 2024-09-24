package errorx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/57blocks/auto-action/cli/internal/pkg/logx"

	"github.com/go-resty/resty/v2"
)

type (
	Errorx struct {
		StatusF     int    `json:"status"`
		CodeF       int    `json:"code"`
		ErrByServer string `json:"error,omitempty"`
		ErrorF      string `json:"message,omitempty"`
	}

	ErrResponse struct {
		Errorx `json:"Error"`
	}
)

func (e *Errorx) Status() int {
	return e.StatusF
}

func (e *Errorx) Code() int {
	return e.CodeF
}

func (e *Errorx) Error() string {
	if len(e.ErrByServer) > 0 {
		return e.ErrByServer
	}
	if len(e.ErrorF) > 0 {
		return e.ErrorF
	}

	return "unknown error"
}

func CatchAndWrap(err error) {
	if err == nil {
		return
	}

	er := new(ErrResponse)
	if errors.Is(err, er) {
		code := http.StatusInternalServerError
		if er.Code() != 0 {
			code = er.Code()
		}
		logx.Logger.Error(
			"server error occurred",
			fmt.Sprintf("error_code_%v", code),
			err.Error(),
		)

		return
	}

	erx := new(Errorx)
	if errors.Is(err, erx) {
		logx.Logger.Error(
			"error occurred",
			fmt.Sprintf("error_code_%v", erx.Code()),
			err.Error(),
		)

		return
	}

	logx.Logger.Error("unrecognized error", "error_msg", err.Error())
	os.Exit(2)
}

func WithRestyResp(resp *resty.Response) error {
	er := new(ErrResponse)
	if err := json.Unmarshal(resp.Body(), er); err != nil {
		logx.Logger.Error("marshal response error", "raw_response", string(resp.Body()))
		return &ErrResponse{
			Errorx: Errorx{
				StatusF: 500,
				CodeF:   500,
				ErrorF:  err.Error(),
			},
		}
	}

	return er
}

func Internal(msg string) error {
	return &Errorx{
		StatusF: 500,
		CodeF:   500,
		ErrorF:  msg,
	}
}

func RestyError(msg string) error {
	return &Errorx{
		StatusF: 0,
		CodeF:   0,
		ErrorF:  fmt.Sprintf("resty error: %s", msg),
	}
}

func BadRequest(msg string) error {
	return &Errorx{
		StatusF: 400,
		CodeF:   400,
		ErrorF:  msg,
	}
}
