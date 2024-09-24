package third_party

import (
	"github.com/57blocks/auto-action/server/internal/third-party/amazonx"
	"github.com/57blocks/auto-action/server/internal/third-party/decrypt"
	"github.com/57blocks/auto-action/server/internal/third-party/jwtx"
	"github.com/57blocks/auto-action/server/internal/third-party/restyx"
)

func Setup() error {
	if err := jwtx.Setup(); err != nil {
		return err
	}
	if err := amazonx.Setup(); err != nil {
		return err
	}
	if err := restyx.Setup(); err != nil {
		return err
	}
	if err := decrypt.Setup(); err != nil {
		return err
	}

	return nil
}
