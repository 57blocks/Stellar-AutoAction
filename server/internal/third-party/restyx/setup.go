package restyx

import "github.com/go-resty/resty/v2"

func Setup() error {
	Conductor = &restyx{client: resty.New()}
	return nil
}
