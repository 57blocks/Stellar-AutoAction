package stellarx

import (
	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/constant"

	"github.com/stellar/go/clients/horizonclient"
)

func Setup() error {
	if config.GlobalConfig.Bound.Name == string(constant.StellarNetworkTypeMainNet) {
		Conductor = &stellar{client: horizonclient.DefaultPublicNetClient}
	} else {
		Conductor = &stellar{client: horizonclient.DefaultTestNetClient}
	}
	return nil
}
