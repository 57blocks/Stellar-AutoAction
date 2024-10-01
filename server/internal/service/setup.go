package service

import (
	"github.com/57blocks/auto-action/server/internal/service/cs"
	"github.com/57blocks/auto-action/server/internal/service/lambda"
	"github.com/57blocks/auto-action/server/internal/service/oauth"
	"github.com/57blocks/auto-action/server/internal/service/wallet"
)

// resource should init after service
func Setup() error {
	cs.NewCubeSignerService()
	lambda.NewLambdaService()
	lambda.NewLambdaResource()
	oauth.NewOAuthService()
	oauth.NewOAuthResource()
	wallet.NewWalletService()

	return nil
}
