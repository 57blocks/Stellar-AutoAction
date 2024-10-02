package wallet

import (
	"context"
	"fmt"
	"strings"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/constant"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/pkg/util"
	"github.com/57blocks/auto-action/server/internal/repo"
	svcCS "github.com/57blocks/auto-action/server/internal/service/cs"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"
	"github.com/57blocks/auto-action/server/internal/third-party/restyx"
	"github.com/57blocks/auto-action/server/internal/third-party/stellarx"

	"github.com/gin-gonic/gin"
	"github.com/stellar/go/clients/horizonclient"
)

//go:generate mockgen -destination ../../testdata/wallet_service_mock.go -package testdata -source service.go Service
type (
	WalletService interface {
		Create(c context.Context) (*dto.RespCreateWallet, error)
		Remove(c context.Context, r *dto.ReqRemoveWallet) error
		List(c context.Context) (*dto.RespListWallets, error)
		Verify(c context.Context, r *dto.ReqVerifyWallet) (*dto.RespVerifyWallet, error)
	}
	service struct {
		oauthRepo repo.OAuth
		csRepo    repo.CubeSigner
		resty     restyx.Resty
		csService svcCS.CSservice
		stellar   stellarx.Stellar
	}
)

var WalletServiceImpl WalletService

func NewWalletService() {
	if WalletServiceImpl == nil {
		repo.NewOAuth()
		repo.NewCubeSigner()

		WalletServiceImpl = &service{
			oauthRepo: repo.OAuthRepo,
			csRepo:    repo.CubeSignerRepo,
			resty:     restyx.Conductor,
			csService: svcCS.CSserviceImpl,
			stellar:   stellarx.Conductor,
		}
	}
}

func (svc *service) Create(c context.Context) (*dto.RespCreateWallet, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get(constant.ClaimIss.Str())
	jwtAccount, _ := ctx.Get(constant.ClaimSub.Str())

	user, err := svc.oauthRepo.FindUserByOrgAcn(c, &dto.ReqOrgAcn{
		OrgName: jwtOrg.(string),
		AcnName: jwtAccount.(string),
	})
	if err != nil {
		return nil, err
	}

	max := config.GlobalConfig.Wallet.Max
	keys, err := svc.csRepo.FindCSKeysByAccount(c, user.ID)
	if err != nil {
		return nil, err
	}
	if len(keys) >= max {
		return nil, errorx.Internal(fmt.Sprintf("the number of wallet address is limited to %d", max))
	}

	csToken, err := svc.csService.CubeSignerToken(c)
	if err != nil {
		return nil, err
	}

	secretName := util.GetSecretName(c, jwtOrg.(string), jwtAccount.(string))
	role, err := svc.csService.GetSecRole(c, secretName)
	if err != nil {
		return nil, err
	}

	keyId, err := svc.resty.AddCSKey(c, csToken)
	if err != nil {
		return nil, err
	}

	if err = svc.resty.AddCSKeyToRole(c, csToken, keyId, role); err != nil {
		return nil, err
	}

	if err := svc.csRepo.SyncCSKey(c, &model.CubeSignerKey{
		AccountID: user.ID,
		Key:       keyId,
		Scopes:    []string{"{sign:blob}"},
	}); err != nil {
		return nil, err
	}

	address, err := util.GetAddressFromCSKey(keyId)
	if err != nil {
		return nil, err
	}

	return &dto.RespCreateWallet{
		Address: address,
	}, nil
}

func (svc *service) Remove(c context.Context, r *dto.ReqRemoveWallet) error {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get(constant.ClaimIss.Str())
	jwtAccount, _ := ctx.Get(constant.ClaimSub.Str())

	user, err := svc.oauthRepo.FindUserByOrgAcn(c, &dto.ReqOrgAcn{
		OrgName: jwtOrg.(string),
		AcnName: jwtAccount.(string),
	})
	if err != nil {
		return err
	}

	keyId := util.GetCSKeyFromAddress(r.Address)
	_, err = svc.csRepo.FindCSKey(c, keyId, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "cube signer key not found") {
			return errorx.Internal(fmt.Sprintf("no existed wallet address found: %s", r.Address))
		}
		return err
	}

	csToken, err := svc.csService.CubeSignerToken(c)
	if err != nil {
		return err
	}

	secretName := util.GetSecretName(c, jwtOrg.(string), jwtAccount.(string))
	role, err := svc.csService.GetSecRole(c, secretName)
	if err != nil {
		return err
	}

	if err = svc.resty.DeleteCSKeyFromRole(c, csToken, keyId, role); err != nil {
		return err
	}

	if err = svc.resty.DeleteCSKey(c, csToken, keyId); err != nil {
		return err
	}

	if err := svc.csRepo.DeleteCSKey(c, keyId, user.ID); err != nil {
		return err
	}

	return nil
}

func (svc *service) List(c context.Context) (*dto.RespListWallets, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get(constant.ClaimIss.Str())
	jwtAccount, _ := ctx.Get(constant.ClaimSub.Str())

	user, err := svc.oauthRepo.FindUserByOrgAcn(c, &dto.ReqOrgAcn{
		OrgName: jwtOrg.(string),
		AcnName: jwtAccount.(string),
	})
	if err != nil {
		return nil, err
	}

	keys, err := svc.csRepo.FindCSKeysByAccount(c, user.ID)
	if err != nil {
		return nil, err
	}

	// convert db data to response result
	response := &dto.RespListWallets{
		Data: make([]dto.RespListWallet, len(keys)),
	}
	for i, key := range keys {
		address, err := util.GetAddressFromCSKey(key.Key)
		if err != nil {
			return nil, err
		}
		response.Data[i] = dto.RespListWallet{
			Address: address,
		}
	}

	return response, nil
}

func (svc *service) Verify(c context.Context, r *dto.ReqVerifyWallet) (*dto.RespVerifyWallet, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get(constant.ClaimIss.Str())
	jwtAccount, _ := ctx.Get(constant.ClaimSub.Str())

	user, err := svc.oauthRepo.FindUserByOrgAcn(c, &dto.ReqOrgAcn{
		OrgName: jwtOrg.(string),
		AcnName: jwtAccount.(string),
	})
	if err != nil {
		return nil, err
	}

	keyId := util.GetCSKeyFromAddress(r.Address)
	_, err = svc.csRepo.FindCSKey(c, keyId, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "cube signer key not found") {
			return nil, errorx.Internal(fmt.Sprintf("no existed wallet address found: %s", r.Address))
		}
		return nil, err
	}

	_, err = svc.stellar.AccountDetail(c, horizonclient.AccountRequest{AccountID: r.Address})
	if err != nil {
		logx.Logger.ERROR(fmt.Sprintf("verify wallet address %s occurred error: %s", r.Address, err.Error()))
		return &dto.RespVerifyWallet{
			Address: r.Address,
			IsValid: false,
		}, nil
	}

	return &dto.RespVerifyWallet{
		Address: r.Address,
		IsValid: true,
	}, nil
}
