package wallet

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/model/cs"
	"github.com/57blocks/auto-action/server/internal/model/oauth"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/pkg/logx"
	"github.com/57blocks/auto-action/server/internal/pkg/restyx"
	csSvc "github.com/57blocks/auto-action/server/internal/service/cs"
	csdto "github.com/57blocks/auto-action/server/internal/service/dto/cs"
	walletdto "github.com/57blocks/auto-action/server/internal/service/dto/wallet"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	Service interface {
		Create(c context.Context, r *http.Request) (*walletdto.RespInfo, error)
	}
	conductor struct{}
)

var (
	Conductor Service
)

func init() {
	if Conductor == nil {
		Conductor = &conductor{}
	}
}

func (cd conductor) Create(c context.Context, r *http.Request) (*walletdto.RespInfo, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get("jwt_organization")
	jwtAccount, _ := ctx.Get("jwt_account")

	// fetch the organization by org from the database
	org := new(oauth.Organization)
	if err := db.Conn(c).Table(org.TableName()).
		Where(map[string]interface{}{
			"name": jwtOrg,
		}).
		First(org).Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("none organization found")
		}

		return nil, errorx.Internal(fmt.Sprintf("find organization by name: %s, occurred error: %s", jwtOrg, err.Error()))
	}

	// fetch the user by account from the database
	user := new(oauth.User)
	if err := db.Conn(c).Table(user.TableName()).
		Where(map[string]interface{}{
			"account":         jwtAccount,
			"organization_id": org.ID,
		}).
		First(user).Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("none account found")
		}
		return nil, errorx.Internal(fmt.Sprintf("find account by name: %s, occurred error: %s", jwtAccount, err.Error()))
	}

	// fetch the organization and role by account and organization from the database
	role := new(cs.CubeSignerRole)
	if err := db.Conn(c).Table(role.TableName()).
		Where(map[string]interface{}{
			"organization_id": org.ID,
			"account_id":      user.ID,
		}).
		First(role).Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("none role found")
		}
		return nil, errorx.Internal(fmt.Sprintf("find role by org_id: %d and account_id: %d, occurred error: %s", org.ID, user.ID, err.Error()))
	}

	csToken, err := csSvc.Conductor.CubeSignerToken(c)
	if err != nil {
		return nil, err
	}

	keyId, err := addCsKey(org, csToken, user)
	if err != nil {
		return nil, err
	}

	err = addKeyToRole(org, role, csToken, keyId)
	if err != nil {
		return nil, err
	}

	err = saveCsKey(c, keyId, role)
	if err != nil {
		return nil, err
	}

	// parse key_id(format: Key#Stellar_<address>) to get the address
	address := strings.Split(keyId, "_")[1]
	return &walletdto.RespInfo{
		Address: address,
	}, nil
}

func addCsKey(org *oauth.Organization, csToken string, user *oauth.User) (string, error) {
	URL := fmt.Sprintf("https://gamma.signer.cubist.dev/v0/org/%s/keys", url.PathEscape(org.CubeSignerOrg))
	var keyResp csdto.KeyResponse
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

func addKeyToRole(org *oauth.Organization, role *cs.CubeSignerRole, csToken string, keyId string) error {
	URL := fmt.Sprintf("https://gamma.signer.cubist.dev/v0/org/%s/roles/%s/add_keys", url.PathEscape(org.CubeSignerOrg), url.PathEscape(role.Role))
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

func saveCsKey(c context.Context, keyId string, role *cs.CubeSignerRole) error {
	csKey := &cs.CubeSignerKey{
		Key:    keyId,
		RoleID: role.ID,
		Scopes: []string{"{sign:blob}"},
	}
	if err := db.Conn(c).
		Table(csKey.TableName()).Create(csKey).
		Error; err != nil {
		return errorx.Internal(err.Error())
	}
	logx.Logger.DEBUG(fmt.Sprintf("save cube signer key to database success: %s", csKey.Key))
	return nil
}
