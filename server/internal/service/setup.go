package service

import (
	"github.com/57blocks/auto-action/server/internal/service/cs"
	"github.com/57blocks/auto-action/server/internal/service/lambda"
	"github.com/57blocks/auto-action/server/internal/service/oauth"
	"github.com/57blocks/auto-action/server/internal/service/wallet"
)

func Setup() error {
	setupServices()

	return nil
}

func setupServices() {
	cs.NewCubeSignerService()
	lambda.NewLambdaService()
	oauth.NewOAuthService()
	wallet.NewWalletService()
}
