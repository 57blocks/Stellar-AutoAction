package organization

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	configx "github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/model"
	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"
	"github.com/57blocks/auto-action/server/internal/pkg/restyx"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/organization"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type (
	Service interface {
		Organization(c context.Context) (*model.Organization, error)
		OrgRoleKey(c context.Context, req *dto.ReqSDKRequired) (*dto.RespOrgRoleKey, error)
		OrgRootSession(c context.Context, req *dto.ReqSDKRequired) (*dto.RespSDKRequired, error)

		OrgSecret(c context.Context) (string, error)

		KeepRootSessionAlive(c context.Context) error
	}
	conductor struct{}
)

var (
	Conductor Service
	awsConfig aws.Config
	smClient  *secretsmanager.Client
)

func init() {
	if Conductor == nil {
		Conductor = &conductor{}
	}
}

func (cd conductor) Organization(c context.Context) (*model.Organization, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errors.New("convert context.Context to gin.Context failed")
	}

	jwtOrg, _ := ctx.Get("jwt_organization")

	org := new(model.Organization)
	if err := db.Conn(c).Table(org.TableName()).
		Where(map[string]interface{}{
			"name": jwtOrg,
		}).
		First(org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("none organization found")
		}

		return nil, errors.Wrap(err, "db error when find organization by name")
	}

	return org, nil
}

func (cd conductor) OrgRoleKey(c context.Context, req *dto.ReqSDKRequired) (*dto.RespOrgRoleKey, error) {
	orkList := make([]*model.CSOrgRoleKey, 0)
	if err := db.Conn(c).
		Table(model.TabNameOrgRoleKeyAbbr()).
		Joins("LEFT JOIN organization AS o ON o.id = ork.organization_id").
		Joins("LEFT JOIN \"user\" AS u ON u.organization_id = o.id").
		Where(map[string]interface{}{
			"o.name": req.Organization,
		}).
		Find(&orkList).Error; err != nil {
		return nil, errors.Wrap(err, "db error when query organization key pairs")
	}

	if len(orkList) == 0 {
		return nil, errors.New("none organization related key pairs found")
	}

	csRoleKeyList := make([]dto.RespCSRoleKey, 0, len(orkList))
	for _, ork := range orkList {
		csRoleKeyList = append(csRoleKeyList, dto.RespCSRoleKey{
			CSRoleID: ork.CSRoleID,
			CSKeyID:  ork.CSKeyID,
			CSScopes: ork.CSScopes,
		})
	}

	return &dto.RespOrgRoleKey{
		CSRoleKeys: csRoleKeyList,
	}, nil
}

func (cd conductor) OrgRootSession(c context.Context, req *dto.ReqSDKRequired) (*dto.RespSDKRequired, error) {
	orgSession := new(model.OrgRootSession)
	if err := db.Conn(c).
		Table(model.TabNameOrgSessionAbbr()).
		Joins("LEFT JOIN organization AS o ON o.id = ors.organization_id").
		Where(map[string]interface{}{
			"o.name": req.Organization,
		}).
		First(orgSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("none organization root session found")
		}

		return nil, errors.Wrap(err, "db error when query organization key pairs")
	}

	org := new(model.Organization)
	if err := db.Conn(c).Table(org.TableName()).
		Where(map[string]interface{}{
			"name": req.Organization,
		}).First(org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("none organization found")
		}

		return nil, errors.Wrap(err, "db error when find organization by name")
	}

	return &dto.RespSDKRequired{
		Token:        orgSession.Token,
		Organization: org.CSOrganizationID,
	}, nil
}

func (cd conductor) OrgSecret(c context.Context) (string, error) {
	var err error

	awsConfig, err := config.LoadDefaultConfig(
		c,
		config.WithRegion(configx.Global.Region),
		config.WithSharedConfigProfile("iamp3ngf3i"), // TODO: only for local
	)

	smClient = secretsmanager.NewFromConfig(awsConfig)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("AutoActionSecretKey-Dev"),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	resp, err := smClient.GetSecretValue(c, input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return "", errors.Wrap(err, "failed to get secret value")
	}

	resMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(*resp.SecretString), &resMap); err != nil {
		return "", errors.Wrap(err, "json unmarshal error when parse secret value")
	}

	return resMap["Server_Secret_Key"].(string), nil
}

// KeepRootSessionAlive
// Login through `cs` manually at the first time, then save the required info to DB,
// Here below will keep the root session alive by calling the refresh endpoint
// with a fixed interval.
//
// The interval is based on the `--refresh-lifetime` in root session initialization:
//
//	`cs login -s google --session-lifetime 31536000 --auth-lifetime 86400 --refresh-lifetime 86400`
//
// The longest session lifetime is 31536000 seconds, which is about 1 year.
//
// Once the root session is expired, login through `cs` again and sync the info to DB.
func (cd conductor) KeepRootSessionAlive(c context.Context) error {
	rootSessions := make([]*model.OrgRootSession, 0)
	if err := db.Conn(c).Table(model.TabNameOrgSession()).
		Where(map[string]interface{}{}).
		Find(rootSessions).Error; err != nil {
		return errors.Wrap(err, "db error when find organization root session")
	}

	if len(rootSessions) == 0 {
		return errors.New("none organization root sessions found")
	}

	results := make([]*dto.RespRefreshRootSession, 0, len(rootSessions))
	for _, session := range rootSessions {
		respRRS := new(dto.RespRefreshRootSession)

		response, err := restyx.Client.R().
			EnableTrace().
			SetHeaders(map[string]string{
				"Content-Type":  "application/json",
				"accept":        "application/json",
				"Authorization": session.Token,
			}).
			SetBody(map[string]interface{}{
				"epoch_num":   session.Epoch,
				"epoch_token": session.EpochToken,
				"other_token": session.EpochRefreshToken,
			}).
			Post(fmt.Sprintf(""))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("resty error: %s\n", err.Error()))
		}
		if response.IsError() {
			return errors.New(fmt.Sprintf("resty response error: %s", string(response.Body())))
		}
		fmt.Println(string(response.Body()))

		if err := json.Unmarshal(response.Body(), respRRS); err != nil {
			return errors.Wrap(err, "json unmarshal error when parse refresh response")
		}

		results = append(results, respRRS)
	}

	go func(c context.Context, results *[]*dto.RespRefreshRootSession) {
		if err := db.Conn(c).Transaction(
			func(tx *gorm.DB) error {
				// TODO: sync results to DB
				return nil
			},
			&sql.TxOptions{
				Isolation: sql.LevelSerializable,
			},
		); err != nil {
			pkgLog.Logger.ERROR(fmt.Sprintf("failed to update organization root session: %s", err.Error()))
		}
	}(c, &results)

	return nil
}
