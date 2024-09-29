package restyx

import (
	"github.com/go-resty/resty/v2"
)

func Setup() error {
	// TODO: remove Client later
	if Client == nil {
		Client = resty.New()
	}
	Conductor = &restyX{client: Client}
	return nil
}
