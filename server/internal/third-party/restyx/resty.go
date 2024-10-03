package restyx

import (
	"context"
	"fmt"
	"net/url"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"github.com/go-resty/resty/v2"
)

//go:generate mockgen -destination ../../testdata/resty_mock.go -package testdata -source resty.go Resty
type (
	Resty interface {
		AddCSRole(c context.Context, csToken string, orgName string, account string) (*dto.RespAddCsRole, error)
		AddCSKey(c context.Context, csToken string) (string, error)
		AddCSKeyToRole(c context.Context, csToken string, keyId string, role string) error
		DeleteCSKey(c context.Context, csToken string, keyId string) error
		DeleteCSKeyFromRole(c context.Context, csToken string, keyId string, role string) error
	}

	restyx struct {
		client *resty.Client
	}
)

var Conductor Resty

func (r *restyx) AddCSRole(c context.Context, csToken string, orgName string, account string) (*dto.RespAddCsRole, error) {
	URL := fmt.Sprintf(
		"%s/v0/org/%s/roles",
		config.GlobalConfig.CS.Endpoint,
		url.PathEscape(config.GlobalConfig.CS.Organization),
	)

	var roleResp dto.RespAddCsRole
	roleName := fmt.Sprintf("%s_%s_Role", orgName, account)
	resp, err := r.client.R().
		SetHeader("Authorization", csToken).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"name": roleName,
		}).
		SetResult(&roleResp).
		Post(URL)
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("create cube signer role occurred error: %s", err.Error()))
	}
	if resp.IsError() {
		return nil, errorx.Internal(fmt.Sprintf("create cube signer role occurred error: %d, %s", resp.StatusCode(), resp.String()))
	}

	logx.Logger.DEBUG(fmt.Sprintf("create cube signer role success: %s", roleName))

	return &roleResp, nil
}

func (r *restyx) AddCSKey(c context.Context, csToken string) (string, error) {
	URL := fmt.Sprintf(
		"%s/v0/org/%s/keys",
		config.GlobalConfig.CS.Endpoint,
		url.PathEscape(config.GlobalConfig.CS.Organization),
	)

	var keyResp dto.RespAddCsKey
	resp, err := r.client.R().
		SetHeader("Authorization", csToken).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"count":    1,
			"key_type": "Ed25519StellarAddr",
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

func (r *restyx) AddCSKeyToRole(c context.Context, csToken string, keyId string, role string) error {
	URL := fmt.Sprintf(
		"%s/v0/org/%s/roles/%s/add_keys",
		config.GlobalConfig.CS.Endpoint,
		url.PathEscape(config.GlobalConfig.CS.Organization),
		url.PathEscape(role),
	)

	resp, err := r.client.R().
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

func (r *restyx) DeleteCSKey(c context.Context, csToken string, keyId string) error {
	URL := fmt.Sprintf(
		"%s/v0/org/%s/keys/%s",
		config.GlobalConfig.CS.Endpoint,
		url.PathEscape(config.GlobalConfig.CS.Organization),
		url.PathEscape(keyId),
	)

	resp, err := r.client.R().
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

func (r *restyx) DeleteCSKeyFromRole(c context.Context, csToken string, keyId string, role string) error {
	URL := fmt.Sprintf(
		"%s/v0/org/%s/roles/%s/keys/%s",
		config.GlobalConfig.CS.Endpoint,
		url.PathEscape(config.GlobalConfig.CS.Organization),
		url.PathEscape(role),
		url.PathEscape(keyId),
	)

	resp, err := r.client.R().
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
