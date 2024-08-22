package restyx

import (
	"github.com/go-resty/resty/v2"
)

var Client *resty.Client

func init() {
	if Client == nil {
		Client = resty.New()
	}
}
