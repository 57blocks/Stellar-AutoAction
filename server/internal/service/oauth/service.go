package oauth

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/pkg/util"
	"github.com/57blocks/auto-action/server/internal/repo"
	svcCS "github.com/57blocks/auto-action/server/internal/service/cs"
	"github.com/57blocks/auto-action/server/internal/third-party/amazonx"
	"github.com/57blocks/auto-action/server/internal/third-party/decrypt"
	"github.com/57blocks/auto-action/server/internal/third-party/jwtx"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"
	"github.com/57blocks/auto-action/server/internal/third-party/restyx"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -destination ../../testdata/oauth_service_mock.go -package testdata -source service.go Service
type (
	OAuthService interface {
		Signup(c context.Context, req dto.ReqSignup) error
		Login(c context.Context, req dto.ReqLogin) (*dto.RespCredential, error)
		Refresh(c context.Context, raw string) (*dto.RespCredential, error)
		Logout(c context.Context, raw string) (*dto.RespLogout, error)
	}
	service struct {
		jwtx      jwtx.JWT
		decrypter decrypt.Decrypter
		oauthRepo repo.OAuth
		amazon    amazonx.Amazon
		resty     restyx.Resty
		csService svcCS.CSservice
	}
)

var OAuthServiceImpl OAuthService

func NewOAuthService() {
	if OAuthServiceImpl == nil {
		repo.NewOAuth()

		OAuthServiceImpl = &service{
			jwtx:      jwtx.RS256,
			decrypter: decrypt.RSADecrypter,
			oauthRepo: repo.OAuthRepo,
			amazon:    amazonx.Conductor,
			resty:     restyx.Conductor,
			csService: svcCS.CSserviceImpl,
		}
	}
}

func (svc *service) Signup(c context.Context, req dto.ReqSignup) error {
	org, err := svc.oauthRepo.FindOrgByName(c, req.Organization)
	if err != nil {
		return err
	}

	user, err := svc.oauthRepo.FindUserByAcn(c, req.Account)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}
	if user != nil {
		return errorx.BadRequest("user already exists")
	}

	csToken, err := svc.csService.CubeSignerToken(c)
	if err != nil {
		return err
	}

	csRole, err := svc.resty.AddCSRole(c, csToken, req.Organization, req.Account)
	if err != nil {
		return err
	}

	awsRole, err := svc.addAwsRole(c, req.Organization, req.Account)
	if err != nil {
		return err
	}

	if err = svc.addAwsSecretKey(c, req.Organization, req.Account, csRole, awsRole); err != nil {
		return err
	}

	rawPwdBytes, err := svc.decrypter.Decrypt([]byte(req.Password))
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPwdBytes), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}
	if err := svc.oauthRepo.CreateUser(c, &model.User{
		OrganizationId: org.ID,
		Account:        req.Account,
		Password:       string(hashedPassword),
		Description:    description,
	}); err != nil {
		return err
	}

	return nil
}

func (svc *service) Login(c context.Context, req dto.ReqLogin) (*dto.RespCredential, error) {
	u, err := svc.oauthRepo.FindUserByOrgAcn(c, &dto.ReqOrgAcn{
		OrgName: req.Organization,
		AcnName: req.Account,
	})
	if err != nil {
		return nil, err
	}

	rawPwdBytes, err := svc.decrypter.Decrypt([]byte(req.Password))
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), rawPwdBytes); err != nil {
		return nil, errorx.BadRequest("password not match")
	}

	// tokens assignment
	now := time.Now().UTC()
	accessID := svc.jwtx.GenerateID()
	accessExp := now.AddDate(0, 0, 7)
	refreshID := svc.jwtx.GenerateID()
	refreshExp := now.AddDate(0, 1, 0)

	access, err := svc.jwtx.Assign(&jwtx.AAClaims{
		StdJWTClaims: jwt.StandardClaims{
			Audience:  config.GlobalConfig.Bound.EndPoint,
			ExpiresAt: accessExp.Unix(),
			Id:        accessID,
			IssuedAt:  now.Unix(),
			Issuer:    req.Organization,
			NotBefore: now.Unix(),
			Subject:   u.Account,
		},
	})
	if err != nil {
		return nil, err
	}

	refresh, err := svc.jwtx.Assign(&jwtx.AAClaims{
		StdJWTClaims: jwt.StandardClaims{
			Audience:  config.GlobalConfig.Bound.EndPoint,
			ExpiresAt: refreshExp.Unix(),
			Id:        refreshID,
			IssuedAt:  now.Unix(),
			Issuer:    req.Organization,
			NotBefore: now.Unix(),
			//NotBefore: accessExp.Unix(), // won't be valid until access token expires
			Subject: u.Account,
		},
	})
	if err != nil {
		return nil, err
	}

	// sync token pairs
	token := &model.Token{
		UserId:         u.ID,
		Access:         access,
		AccessID:       accessID,
		AccessExpires:  accessExp,
		Refresh:        refresh,
		RefreshID:      refreshID,
		RefreshExpires: refreshExp,
	}
	if err := svc.oauthRepo.SyncToken(c, token); err != nil {
		return nil, err
	}

	// build response
	resp := dto.BuildRespCred(
		dto.WithAccount(req.Account),
		dto.WithOrganization(req.Organization),
		dto.WithNetwork(config.GlobalConfig.Bound.Name),
		dto.WithTokenPair(jwtx.TokenPair{
			Access:  access,
			Refresh: refresh,
		}),
	)

	return resp, nil
}

func (svc *service) Refresh(c context.Context, raw string) (*dto.RespCredential, error) {
	jwtClaims, err := svc.jwtx.Parse(raw)
	if err != nil {
		return nil, err
	}

	aaClaims, ok := jwtClaims.(*jwtx.AAClaims)
	if !ok {
		return nil, errorx.Internal("failed to assert JWT claims as AAClaims")
	}

	token, err := svc.oauthRepo.FindTokenByRefreshID(c, aaClaims.StdJWTClaims.Id) // use refresh id
	if err != nil {
		return nil, errorx.UnauthorizedWithMsg("invalid refresh token")
	}

	now := time.Now().UTC()
	accessID := svc.jwtx.GenerateID()
	accessExp := now.AddDate(0, 0, 7)

	access, err := svc.jwtx.Assign(&jwtx.AAClaims{
		StdJWTClaims: jwt.StandardClaims{
			Audience:  config.GlobalConfig.Bound.EndPoint,
			ExpiresAt: accessExp.Unix(),
			Id:        accessID,
			IssuedAt:  now.Unix(),
			Issuer:    aaClaims.StdJWTClaims.Issuer,
			NotBefore: now.Unix(),
			Subject:   aaClaims.StdJWTClaims.Subject,
		},
	})
	if err != nil {
		return nil, err
	}

	// save tokens association
	token.Access = access
	token.AccessID = accessID
	token.AccessExpires = accessExp
	token.UpdatedAt = &now

	if err := svc.oauthRepo.SyncToken(c, token); err != nil {
		return nil, err
	}

	resp := dto.BuildRespCred(
		dto.WithAccount(aaClaims.StdJWTClaims.Subject),
		dto.WithOrganization(aaClaims.StdJWTClaims.Issuer),
		dto.WithNetwork(config.GlobalConfig.Bound.Name),
		dto.WithTokenPair(jwtx.TokenPair{
			Access:  access,
			Refresh: raw,
		}),
	)

	return resp, nil
}

func (svc *service) Logout(c context.Context, raw string) (*dto.RespLogout, error) {
	if err := svc.oauthRepo.DeleteTokenByAccess(c, raw); err != nil {
		return nil, err
	}

	return new(dto.RespLogout), nil
}

func (svc *service) addAwsRole(c context.Context, orgName string, account string) (*iam.CreateRoleOutput, error) {
	roleName := util.GetRoleName(c, orgName, account)
	assumeRolePolicyDocument := `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": {
				"Service": [
					"events.amazonaws.com",
					"lambda.amazonaws.com",
					"scheduler.amazonaws.com"
				]
			},
			"Action": "sts:AssumeRole"
		}
	]
}`
	createRoleInput := &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(assumeRolePolicyDocument),
		Description:              aws.String(fmt.Sprintf("Role for AutoAction User %s-%s", orgName, account)),
	}
	role, err := svc.amazon.CreateRole(c, createRoleInput)
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("create aws role occurred error: %s", err.Error()))
	}

	policyDocument := fmt.Sprintf(`{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": [
				"lambda:InvokeFunction",
				"lambda:GetFunction"
			],
			"Resource": "arn:aws:lambda:*:*:function:%s-%s*"
		},
		{
			"Effect": "Allow",
			"Action": [
				"secretsmanager:GetSecretValue",
				"secretsmanager:DescribeSecret"
			],
			"Resource": "*"
		},
		{
			"Effect": "Allow",
			"Action": [
				"logs:CreateLogGroup",
				"logs:CreateLogStream",
				"logs:PutLogEvents"
			],
			"Resource": "*"
		}
	]
}`,
		orgName,
		account,
	)
	putRolePolicyInput := &iam.PutRolePolicyInput{
		PolicyName:     aws.String(fmt.Sprintf("%s-policy", roleName)),
		PolicyDocument: aws.String(policyDocument),
		RoleName:       aws.String(roleName),
	}
	_, err = svc.amazon.PutRolePolicy(c, putRolePolicyInput)
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("put role policy for aws role occurred error: %s", err.Error()))
	}

	logx.Logger.DEBUG(fmt.Sprintf("create aws role success: %s", *role.Role.Arn))

	return role, nil
}

func (svc *service) addAwsSecretKey(
	c context.Context,
	orgName string,
	account string,
	csRole *dto.RespAddCsRole,
	awsRole *iam.CreateRoleOutput,
) error {
	secretName := util.GetSecretName(c, orgName, account)
	secretValue := fmt.Sprintf(`{"cs_role": "%s"}`, csRole.RoleId)
	input := &secretsmanager.CreateSecretInput{
		Name:         aws.String(secretName),
		SecretString: aws.String(secretValue),
		Description:  aws.String(fmt.Sprintf("Secret for %s-%s", orgName, account)),
	}
	awsSecretKey, err := svc.amazon.CreateSecret(c, input)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("create aws secret key occurred error: %s", err.Error()))
	}

	// sleep times to make sure the secret is created, default 10 seconds
	sleepTime, _ := strconv.Atoi(config.GlobalConfig.Amazon.SecretCreateSleepTime)
	time.Sleep(time.Duration(sleepTime) * time.Second)

	resourcePolicy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"AWS": [ "%s", "%s" ]
				},
				"Action": [ "secretsmanager:GetSecretValue", "secretsmanager:DescribeSecret" ],
				"Resource": "%s"
			},
			{
				"Effect": "Deny",
				"Principal" : "*",
				"Action" : "secretsmanager:GetSecretValue",
				"Resource" : "%s",
				"Condition": {
					"StringNotEquals": {
						"aws:PrincipalArn": [ "%s", "%s" ]
					}
				}
			}
		]
	}`,
		config.GlobalConfig.Amazon.EcsTaskRole,
		*awsRole.Role.Arn,
		*awsSecretKey.ARN,
		*awsSecretKey.ARN,
		config.GlobalConfig.Amazon.EcsTaskRole,
		*awsRole.Role.Arn,
	)
	policyInput := &secretsmanager.PutResourcePolicyInput{
		SecretId:       aws.String(secretName),
		ResourcePolicy: aws.String(resourcePolicy),
	}
	_, err = svc.amazon.PutResourcePolicy(c, policyInput)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("put resource policy for aws secret key occurred error: %s", err.Error()))
	}

	logx.Logger.DEBUG(fmt.Sprintf("create aws secret key success: %s", *awsSecretKey.ARN))

	return nil
}
