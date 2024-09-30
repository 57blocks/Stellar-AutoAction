package stellarx

import (
	"context"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
)

//go:generate mockgen -destination ./stellar_mock.go -package stellarx -source stellar.go Stellar
type (
	Stellar interface {
		AccountDetail(c context.Context, req horizonclient.AccountRequest) (horizon.Account, error)
	}

	stellar struct {
		client *horizonclient.Client
	}
)

var Conductor Stellar

func (s *stellar) AccountDetail(c context.Context, req horizonclient.AccountRequest) (horizon.Account, error) {
	return s.client.AccountDetail(req)
}
