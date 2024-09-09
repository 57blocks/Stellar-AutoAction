package cs

import (
	"context"
	"encoding/json"
	"fmt"

	configx "github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/repo"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type (
	Service interface {
		APIKey(c context.Context) (string, error)
		ToSign(c context.Context, req *dto.ReqToSign) ([]*dto.RespToSign, error)
		CubeSignerToken(c context.Context) (string, error)
	}
	conductor struct {
		csRepo repo.CubeSigner
	}
)

var (
	Conductor Service
	awsConfig aws.Config
	smClient  *secretsmanager.Client
)

func NewCubeSignerService() {
	if Conductor == nil {
		repo.NewCubeSigner()

		Conductor = &conductor{
			csRepo: repo.CubeSignerImpl,
		}
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

func (cd conductor) ToSign(c context.Context, req *dto.ReqToSign) ([]*dto.RespToSign, error) {
	toSigns, err := cd.csRepo.ToSign(c, &dto.ReqToSign{
		Organization: req.Organization,
		Account:      req.Account,
	})
	if err != nil {
		return nil, err
	}

	resp := make([]*dto.RespToSign, 0, len(toSigns))
	for _, toSign := range toSigns {
		if len(toSign.Keys) == 0 {
			continue
		}

		ks := make([]dto.RespCSKey, 0, len(toSign.Keys))
		for _, k := range toSign.Keys {
			ks = append(ks, dto.RespCSKey{
				Key:    k.Key,
				Scopes: k.Scopes,
			})
		}
		resp = append(resp, &dto.RespToSign{
			Organization: toSign.Organization,
			Account:      toSign.Account,
			Role:         toSign.Role,
			Keys:         ks,
		})
	}

	return resp, nil
}

func (cd conductor) CubeSignerToken(c context.Context) (string, error) {
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
		SecretId:     aws.String("AA_CS_Token"),
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

	return resMap["token"].(string), nil
}
