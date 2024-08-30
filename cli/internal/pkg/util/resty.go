package util

import (
	"encoding/json"
	"github.com/pkg/errors"

	"github.com/go-resty/resty/v2"
)

type Error struct {
	Message string `json:"message"`
	Notes   string `json:"notes"`
	Status  int    `json:"status"`
}

func (e Error) Error() string {
	return e.Message
}

func HasError(resp *resty.Response) error {
	var err Error
	if err := json.Unmarshal(resp.Body(), &err); err != nil {
		return errors.New(err.Error())
	}

	if err.Status == 0 || (resp.StatusCode() > 199 && resp.StatusCode() < 300) {
		return nil
	}

	return err
}
