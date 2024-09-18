package cs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/repo"
	"github.com/57blocks/auto-action/server/internal/third-party/amazonx"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type (
	Service interface {
		APIKey(c context.Context) (string, error)
		ToSign(c context.Context, req *dto.ReqToSign) (*dto.RespCSKey, error)
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

func (svc *service) APIKey(c context.Context) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("AA_API_KEY"),
		VersionStage: aws.String("AWSCURRENT"),
	}

	secretResp, err := svc.amazon.GetSecretValue(c, input)
	if err != nil {
		return "", errorx.Internal(err.Error())
	}

	resMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(*secretResp.SecretString), &resMap); err != nil {
		return "", errorx.Internal(fmt.Sprintf("json unmarshal error when parse secret value: %s", err.Error()))
	}

	return resMap["api_key"].(string), nil
}

func (svc *service) ToSign(c context.Context, req *dto.ReqToSign) (*dto.RespCSKey, error) {
	user, err := svc.oauthRepo.FindUserByOrgAcn(c, &dto.ReqOrgAcn{
		OrgName: req.Organization,
		AcnName: req.Account,
	})
	if err != nil {
		return nil, err
	}

	// Key sample: Key#Stellar_ABCDEFG
	forSign, err := svc.csRepo.ToSign(c, user.ID, fmt.Sprintf("Key#Stellar_%s", req.From))
	if err != nil {
		return nil, err
	}

	forSign.Organization = config.GlobalConfig.CS.Organization

	return forSign, nil
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
