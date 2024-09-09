package wallet

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/repo"
	csSvc "github.com/57blocks/auto-action/server/internal/service/cs"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"
	"github.com/57blocks/auto-action/server/internal/third-party/restyx"

	"github.com/gin-gonic/gin"
)

type (
	Service interface {
		Create(c context.Context, r *http.Request) (*dto.CreateWalletRespInfo, error)
	}
	conductor struct {
		oauthRepo repo.OAuth
		csRepo    repo.CubeSigner
	}
)

var (
	Conductor Service
)

func NewWalletService() {
	if Conductor == nil {
		repo.NewOAuth()
		repo.NewCubeSigner()

		Conductor = &conductor{
			oauthRepo: repo.OAuthImpl,
			csRepo:    repo.CubeSignerImpl,
		}
	}
}

func (cd *conductor) Create(c context.Context, r *http.Request) (*dto.CreateWalletRespInfo, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get("jwt_organization")
	jwtAccount, _ := ctx.Get("jwt_account")

	// fetch the organization by org from the database
	org, err := cd.oauthRepo.FindOrgByName(c, jwtOrg.(string))
	if err != nil {
		return nil, err
	}

	// fetch the user by account from the database
	//user := new(model.User)
	user, err := cd.oauthRepo.FindUserByOrgAcn(c, dto.ReqOrgAcn{
		OrgName: org.Name,
		AcnName: jwtAccount.(string),
	})
	if err != nil {
		return nil, err
	}

	// fetch the organization and role by account and organization from the database
	role, err := cd.csRepo.FindCSByOrgAcn(c, &dto.ReqCSRole{
		OrgID: org.ID,
		AcnID: user.ID,
	})
	if err != nil {
		return nil, err
	}

	csToken, err := csSvc.Conductor.CubeSignerToken(c)
	if err != nil {
		return nil, err
	}

	keyId, err := addCSKey(org, csToken, user)
	if err != nil {
		return nil, err
	}

	err = addKeyToRole(org, role, csToken, keyId)
	if err != nil {
		return nil, err
	}

	err = cd.saveCSKey(c, keyId, role)
	if err != nil {
		return nil, err
	}

	// parse key_id(format: Key#Stellar_<address>) to get the address
	address := strings.Split(keyId, "_")[1]

	return &dto.CreateWalletRespInfo{
		Address: address,
	}, nil
}

func addCSKey(org *dto.RespOrg, csToken string, user *dto.RespUser) (string, error) {
	URL := fmt.Sprintf("%s/v0/org/%s/keys", config.GlobalConfig.CS.Endpoint, url.PathEscape(org.CubeSignerOrg))
	var keyResp dto.AddCsKeyResponse
	resp, err := restyx.Client.R().
		SetHeader("Authorization", csToken).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"count":    1,
			"key_type": "Ed25519StellarAddr",
			"owner":    user.UserKey,
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

func addKeyToRole(org *dto.RespOrg, role *dto.RespCSRole, csToken string, keyId string) error {
	URL := fmt.Sprintf("%s/v0/org/%s/roles/%s/add_keys", config.GlobalConfig.CS.Endpoint, url.PathEscape(org.CubeSignerOrg), url.PathEscape(role.Role))
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

func (cd *conductor) saveCSKey(c context.Context, keyId string, role *dto.RespCSRole) error {
	csKey := &model.CubeSignerKey{
		Key:    keyId,
		RoleID: role.ID,
		Scopes: []string{"{sign:blob}"},
	}
	if err := cd.csRepo.SyncCSKey(c, csKey); err != nil {
		return err
	}
	logx.Logger.DEBUG(fmt.Sprintf("save cube signer key to database success: %s", csKey.Key))

	return nil
}
