package cs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	configx "github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/cs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"gorm.io/gorm"
)

type (
	Service interface {
		APIKey(c context.Context) (string, error)
		ToSign(c context.Context, req *dto.ReqToSign) (*dto.RespToSign, error)
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

func (cd conductor) APIKey(c context.Context) (string, error) {
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
		VersionStage: aws.String("AWSCURRENT"),
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

func (cd conductor) ToSign(c context.Context, req *dto.ReqToSign) (*dto.RespToSign, error) {
	org := new(model.Organization)
	if err := db.Conn(c).Table(org.TableName()).
		Where(map[string]interface{}{
			"name": req.Organization,
		}).
		First(org).Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("none organization found to sign")
		}

		return nil, errorx.Internal(fmt.Sprintf("db error: %s", err.Error()))
	}

	account := new(model.User)
	if err := db.Conn(c).Table(account.TableName()).
		Where(map[string]interface{}{
			"account": req.Account,
		}).
		First(account).Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("none account found to sign")
		}

		return nil, errorx.Internal(fmt.Sprintf("db error: %s", err.Error()))
	}

	role := new(dto.Role)
	if err := db.Conn(c).
		Table(model.TabNameCSRole()).
		//Preload("Keys", func(db *gorm.DB) *gorm.DB {
		//	return db.Table(model.TabNameCSKey())
		//}).
		Where(map[string]interface{}{
			"organization_id": org.ID,
			"account_id":      account.ID,
		}).
		First(role).Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("none role found to sign")
		}

		return nil, errorx.Internal(err.Error())
	}

	resp := &dto.RespToSign{
		Organization: org.CubeSignerOrg,
		Role:         role.Role,
		//Keys:         role.Keys,
	}

	return resp, nil
}
