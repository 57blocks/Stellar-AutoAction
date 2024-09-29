package restyx

import (
	"fmt"
	"net/url"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"github.com/go-resty/resty/v2"
)

// TODO: remove Client later
var Client *resty.Client

//go:generate mockgen -destination ./resty_mock.go -package restyx -source resty.go Resty
type (
	Resty interface {
		AddCSRole(csToken string, orgName string, account string) (*dto.RespAddCsRole, error)
	}

	restyX struct {
		client *resty.Client
	}
)

var Conductor Resty

func (r *restyX) AddCSRole(csToken string, orgName string, account string) (*dto.RespAddCsRole, error) {
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
