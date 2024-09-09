package restyx

import "github.com/go-resty/resty/v2"

var Client *resty.Client

func Setup() error {
	if Client == nil {
		Client = resty.New()
	}

	return nil
}
