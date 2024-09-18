package cs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/repo"
	"github.com/57blocks/auto-action/server/internal/third-party/amazonx"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type (
	Service interface {
		CubeSignerToken(c context.Context) (string, error)
		GetSecRole(c context.Context, secret string) (string, error)
	}
	service struct {
		csRepo    repo.CubeSigner
		oauthRepo repo.OAuth
		amazon    amazonx.Amazon
	}
)

func NewCubeSignerService() {
	if ServiceImpl == nil {
		repo.NewCubeSigner()
		repo.NewOAuth()

		ServiceImpl = &service{
			csRepo:    repo.CubeSignerRepo,
			oauthRepo: repo.OAuthRepo,
			amazon:    amazonx.Conductor,
		}
	}
}

func (svc *service) CubeSignerToken(c context.Context) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("AA_CS_Token"),
		VersionStage: aws.String("AWSCURRENT"),
	}

	resp, err := svc.amazon.GetSecretValue(c, input)
	if err != nil {
		return "", errorx.Internal(err.Error())
	}

	resMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(*resp.SecretString), &resMap); err != nil {
		return "", errorx.Internal(fmt.Sprintf("json unmarshal error when parse secret value: %s", err.Error()))
	}

	return resMap["token"].(string), nil
}

func (svc *service) GetSecRole(c context.Context, secret string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secret),
		VersionStage: aws.String("AWSCURRENT"),
	}

	resp, err := svc.amazon.GetSecretValue(c, input)
	if err != nil {
		return "", errorx.Internal(err.Error())
	}

	resMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(*resp.SecretString), &resMap); err != nil {
		return "", errorx.Internal(fmt.Sprintf("json unmarshal error when parse secret value: %s", err.Error()))
	}

	return resMap["cs_role"].(string), nil
}
