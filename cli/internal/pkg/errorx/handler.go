package errorx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"os"

	"github.com/go-resty/resty/v2"
)

type (
	Errorx struct {
		StatusF   int    `json:"status"`
		CodeF     int    `json:"code"`
		ErrorF    string `json:"error,omitempty"`
		ErrorMsgF string `json:"message,omitempty"`
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
	return e.ErrorMsgF
}

func CatchAndWrap(err error) {
	if err == nil {
		return
	}

	er := new(ErrResponse)
	if errors.As(err, &er) {
		logx.Logger.Error("error occurred", fmt.Sprintf("error_code_%v", er.Code()), err.Error())

		return
	}

	erx := new(Errorx)
	if errors.As(err, &erx) {
		logx.Logger.Error("error occurred", fmt.Sprintf("error_code_%v", erx.Code()), err.Error())

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
				StatusF:   500,
				CodeF:     500,
				ErrorF:    err.Error(),
				ErrorMsgF: "json unmarshal failed",
			},
		}
	}

	return er
}

func Internal(msg string) error {
	return &Errorx{
		StatusF:   500,
		CodeF:     500,
		ErrorMsgF: msg,
	}
}

func RestyError(msg string) error {
	return &Errorx{
		StatusF:   0,
		CodeF:     0,
		ErrorMsgF: fmt.Sprintf("resty error: %s", msg),
	}
}

func BadRequest(msg string) error {
	return &Errorx{
		StatusF:   400,
		CodeF:     400,
		ErrorMsgF: msg,
	}
}
