package organization

import (
	"context"
	"encoding/json"
	configx "github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/model"
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
		OrgRoleKey(c context.Context, req *dto.ReqKeys) (*dto.RespOrgRoleKey, error)
		//OrgRootSession(c context.Context, req *dto.ReqKeys) (*dto.RespKeys, error)

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

func (cd conductor) OrgRoleKey(c context.Context, req *dto.ReqKeys) (*dto.RespOrgRoleKey, error) {
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
