package organization

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	configx "github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/organization"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	Service interface {
		Organization(c context.Context) (*model.Organization, error)
		OrgRoleKey(c context.Context, req *dto.ReqKeys) (*dto.RespOrgRoleKey, error)
		OrgSecret(c context.Context) (string, error)
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
		return nil, errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get("jwt_organization")

	org := new(model.Organization)
	if err := db.Conn(c).Table(org.TableName()).
		Where(map[string]interface{}{
			"name": jwtOrg,
		}).
		First(org).Error; err != nil {
		if errors.As(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("none organization found")
		}

		return nil, errorx.Internal(fmt.Sprintf("find organization by name: %s, occurred error: %s", jwtOrg, err.Error()))
	}

	return org, nil
}

func (cd conductor) OrgRoleKey(c context.Context, req *dto.ReqKeys) (*dto.RespOrgRoleKey, error) {
	orkList := make([]*model.CSOrgRoleKey, 0)
	if err := db.Conn(c).
		Table(model.TabNameOrgRoleKeyAbbr()).
		Joins("LEFT JOIN organization AS o ON o.id = ork.organization_id").
		Joins("LEFT JOIN \"user\" AS u ON u.organization_id = o.id").
		Where(map[string]interface{}{
			"o.name":    req.Organization,
			"u.account": req.Account,
		}).
		Find(&orkList).Error; err != nil {
		return nil, errorx.Internal(err.Error())

	}

	if len(orkList) == 0 {
		return nil, errorx.NotFound("none organization related key pairs found")
	}

	org := new(model.Organization)
	if err := db.Conn(c).Table(org.TableName()).
		Where(map[string]interface{}{
			"id": orkList[0].ID,
		}).
		First(org).Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("none organization found")
		}

		return nil, errorx.Internal(fmt.Sprintf("db error: %s", err.Error()))
	}

	csRoleKeyList := make([]dto.RespCSRoleKey, 0, len(orkList))
	for _, ork := range orkList {
		csRoleKeyList = append(csRoleKeyList, dto.RespCSRoleKey{
			CSOrgID:  org.CSOrganizationID,
			CSRoleID: ork.CSRoleID,
			CSKeyID:  ork.CSKeyID,
			CSScopes: ork.CSScopes,
		})
	}

	return &dto.RespOrgRoleKey{
		CSRoleKeys: csRoleKeyList,
	}, nil
}

func (cd conductor) OrgSecret(c context.Context) (string, error) {
	var err error

	awsConfig, err = config.LoadDefaultConfig(
		c,
		config.WithRegion(configx.GlobalConfig.Region),
	)
	if err != nil {
		return "", errorx.AmazonConfig(err.Error())
	}

	smClient = secretsmanager.NewFromConfig(awsConfig)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("AA_API_KEY"),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	resp, err := smClient.GetSecretValue(c, input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return "", errorx.Internal(err.Error())
	}

	resMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(*resp.SecretString), &resMap); err != nil {
		return "", errorx.Internal(fmt.Sprintf("json unmarshal error when parse secret value: %s", err.Error()))
	}

	return resMap["api_key"].(string), nil
}
