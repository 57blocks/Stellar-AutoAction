package wallet

import (
	"context"
	"fmt"
	"net/url"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/pkg/util"
	"github.com/57blocks/auto-action/server/internal/repo"
	svcCS "github.com/57blocks/auto-action/server/internal/service/cs"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"
	"github.com/57blocks/auto-action/server/internal/third-party/restyx"

	"github.com/gin-gonic/gin"
)

type (
	Service interface {
		Create(c context.Context) (*dto.CreateWalletRespInfo, error)
		Remove(c context.Context, r *dto.RemoveWalletReqInfo) error
		List(c context.Context) (*dto.ListWalletsResponse, error)
	}
	service struct {
		oauthRepo repo.OAuth
		csRepo    repo.CubeSigner
	}
)

func NewWalletService() {
	if ServiceImpl == nil {
		repo.NewOAuth()
		repo.NewCubeSigner()

		ServiceImpl = &service{
			oauthRepo: repo.OAuthRepo,
			csRepo:    repo.CubeSignerRepo,
		}
	}
}

func (svc *service) Create(c context.Context) (*dto.CreateWalletRespInfo, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get("jwt_organization")
	jwtAccount, _ := ctx.Get("jwt_account")

	user, err := svc.oauthRepo.FindUserByOrgAcn(c, &dto.ReqOrgAcn{
		OrgName: jwtOrg.(string),
		AcnName: jwtAccount.(string),
	})
	if err != nil {
		return nil, err
	}

	csToken, err := svcCS.ServiceImpl.CubeSignerToken(c)
	if err != nil {
		return nil, err
	}

	keyId, err := svc.addCSKey(csToken, user)
	if err != nil {
		return nil, err
	}

	if err = svc.addKeyToRole(csToken, keyId); err != nil {
		return nil, err
	}

	if err = svc.saveCSKey(c, user.ID, keyId); err != nil {
		return nil, err
	}

	return &dto.CreateWalletRespInfo{
		Address: util.GetAddressFromCSKey(keyId),
	}, nil
}

func (svc *service) Remove(c context.Context, r *dto.RemoveWalletReqInfo) error {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get("jwt_organization")
	jwtAccount, _ := ctx.Get("jwt_account")

	user, err := svc.oauthRepo.FindUserByOrgAcn(c, &dto.ReqOrgAcn{
		OrgName: jwtOrg.(string),
		AcnName: jwtAccount.(string),
	})
	if err != nil {
		return err
	}

	keyId := fmt.Sprintf("Key#Stellar_%s", r.Address)
	csKey, err := svc.csRepo.FindCSKey(c, keyId, user.ID)
	if err != nil {
		return err
	}
	if csKey == nil {
		return errorx.Internal(fmt.Sprintf("no existed wallet address found: %s", r.Address))
	}

	csToken, err := svcCS.ServiceImpl.CubeSignerToken(c)
	if err != nil {
		return err
	}

	if err = svc.deleteKeyFromRole(csToken, keyId); err != nil {
		return err
	}

	if err = svc.deleteCSKey(csToken, keyId); err != nil {
		return err
	}

	if err := svc.csRepo.DeleteCSKey(c, keyId, user.ID); err != nil {
		return err
	}

	return nil
}

func (svc *service) List(c context.Context) (*dto.ListWalletsResponse, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get("jwt_organization")
	jwtAccount, _ := ctx.Get("jwt_account")

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
	response := &dto.ListWalletsResponse{
		Data: make([]dto.ListWalletRespInfo, len(keys)),
	}
	for i, key := range keys {
		response.Data[i] = dto.ListWalletRespInfo{
			Address: util.GetAddressFromCSKey(key.Key),
		}
	}

	return response, nil
}

func (svc *service) addCSKey(csToken string, user *dto.RespUser) (string, error) {
	URL := fmt.Sprintf(
		"%s/v0/org/%s/keys",
		config.GlobalConfig.CS.Endpoint,
		url.PathEscape(config.GlobalConfig.CS.Organization),
	)

	// TODO: using member to do the request for ut
	var keyResp dto.AddCsKeyResponse
	resp, err := restyx.Client.R().
		SetHeader("Authorization", csToken).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"count":    1,
			"key_type": "Ed25519StellarAddr",
			"owner":    user.CubeSignerUser,
			"policy":   []string{"AllowRawBlobSigning"},
		}).
		SetResult(&keyResp).
		Post(URL)
	if err != nil {
		return "", errorx.Internal(fmt.Sprintf("create cube signer key occurred error: %s", err.Error()))
	}
	if resp.IsError() {
		return "", errorx.Internal(fmt.Sprintf("create cube signer key occurred error: %d, %s", resp.StatusCode(), resp.String()))
	}

	keyId := keyResp.Keys[0].KeyID
	logx.Logger.DEBUG(fmt.Sprintf("create cube signer key success: %s", keyId))

	return keyId, nil
}

func (svc *service) addKeyToRole(csToken string, keyId string) error {
	URL := fmt.Sprintf(
		"%s/v0/org/%s/roles/%s/add_keys",
		config.GlobalConfig.CS.Endpoint,
		url.PathEscape(config.GlobalConfig.CS.Organization),
		url.PathEscape(config.GlobalConfig.CS.Role),
	)

	// TODO: using member to do the request for ut
	resp, err := restyx.Client.R().
		SetHeader("Authorization", csToken).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"key_ids": []string{keyId},
		}).
		Put(URL)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("add cube signer key to role occurred error: %s", err.Error()))
	}
	if resp.IsError() {
		return errorx.Internal(fmt.Sprintf("add cube signer key to role occurred error: %d, %s", resp.StatusCode(), resp.String()))
	}

	logx.Logger.DEBUG(fmt.Sprintf("add cube signer key to role success: %s", resp.String()))

	return nil
}

func (svc *service) saveCSKey(c context.Context, userID uint64, key string) error {
	csKey := &model.CubeSignerKey{
		AccountID: userID,
		Key:       key,
		Scopes:    []string{"{sign:blob}"},
	}
	if err := svc.csRepo.SyncCSKey(c, csKey); err != nil {
		return err
	}

	logx.Logger.DEBUG(fmt.Sprintf("save cube signer key to database success: %s", csKey.Key))

	return nil
}

func (svc *service) deleteCSKey(csToken string, keyId string) error {
	URL := fmt.Sprintf(
		"%s/v0/org/%s/keys/%s",
		config.GlobalConfig.CS.Endpoint,
		url.PathEscape(config.GlobalConfig.CS.Organization),
		url.PathEscape(keyId),
	)

	// TODO: using member to do the request for ut
	resp, err := restyx.Client.R().
		SetHeader("Authorization", csToken).
		SetHeader("Content-Type", "application/json").
		Delete(URL)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("delete cube signer key occurred error: %s", err.Error()))
	}
	if resp.IsError() {
		return errorx.Internal(fmt.Sprintf("delete cube signer key occurred error: %d, %s", resp.StatusCode(), resp.String()))
	}
	logx.Logger.DEBUG(fmt.Sprintf("delete cube signer key success: %s", keyId))

	return nil
}

func (svc *service) deleteKeyFromRole(csToken string, keyId string) error {
	URL := fmt.Sprintf(
		"%s/v0/org/%s/roles/%s/keys/%s",
		config.GlobalConfig.CS.Endpoint,
		url.PathEscape(config.GlobalConfig.CS.Organization),
		url.PathEscape(config.GlobalConfig.CS.Role),
		url.PathEscape(keyId),
	)
	logx.Logger.INFO(fmt.Sprintf("delete cube signer key from role URL: %s", URL))
	resp, err := restyx.Client.R().
		SetHeader("Authorization", csToken).
		SetHeader("Content-Type", "application/json").
		Delete(URL)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("delete cube signer key from role occurred error: %s", err.Error()))
	}
	if resp.IsError() {
		return errorx.Internal(fmt.Sprintf("delete cube signer key from role occurred error: %d, %s", resp.StatusCode(), resp.String()))
	}
	logx.Logger.DEBUG(fmt.Sprintf("delete cube signer key from role success: %s", resp.String()))

	return nil
}
