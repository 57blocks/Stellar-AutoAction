package restyx

import "github.com/go-resty/resty/v2"

func Setup() error {
	Conductor = &restyX{client: resty.New()}
	return nil
}
