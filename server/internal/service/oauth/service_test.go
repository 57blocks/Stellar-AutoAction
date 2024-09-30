package oauth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/repo"
	svcCS "github.com/57blocks/auto-action/server/internal/service/cs"
	"github.com/57blocks/auto-action/server/internal/third-party/amazonx"
	"github.com/57blocks/auto-action/server/internal/third-party/decrypt"
	"github.com/57blocks/auto-action/server/internal/third-party/jwtx"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"
	"github.com/57blocks/auto-action/server/internal/third-party/restyx"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Before test, setup log and config
func TestMain(m *testing.M) {
	os.Setenv("BOUND_ENDPOINT", "http://localhost:8080")
	os.Setenv("BOUND_NAME", "Horizon-Testnet")
	os.Setenv("AWS_SECRET_CREATE_SLEEP_TIME", "0")
	config.Setup("../../config/")
	testConfig := config.Configuration{
		Log: config.Log{
			Level:    "debug",
			Encoding: "json",
		},
	}
	logx.Setup(&testConfig)

	os.Exit(m.Run())
}

func TestLogoutSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	expectedRaw := "test-token"
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockOAuthRepo.EXPECT().DeleteTokenByAccess(ctx, expectedRaw).
		Return(nil)

	svc := &service{
		oauthRepo: mockOAuthRepo,
	}

	resp, err := svc.Logout(ctx, expectedRaw)
	assert.NoError(t, err)
	assert.Equal(t, new(dto.RespLogout), resp)
}

func TestLogoutFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	expectedRaw := "test-token"
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockOAuthRepo.EXPECT().DeleteTokenByAccess(ctx, expectedRaw).
		Return(errors.New("failed to delete token"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
	}

	resp, err := svc.Logout(ctx, expectedRaw)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestRefreshSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockJWT := jwtx.NewMockJWT(ctrl)

	raw := "test-token"
	jwtClaimId := "test-id"
	jwtClaimSubject := "account_name"
	jwtClaimIssuer := "org_name"
	newAccessToken := "new-access-token"
	accessID := "1"

	mockJWT.EXPECT().Parse(raw).Return(&jwtx.AAClaims{
		StdJWTClaims: jwt.StandardClaims{
			Id:      jwtClaimId,
			Subject: jwtClaimSubject,
			Issuer:  jwtClaimIssuer,
		},
	}, nil)

	mockOAuthRepo.EXPECT().FindTokenByRefreshID(ctx, jwtClaimId).
		Return(&model.Token{
			UserId: 1,
		}, nil)

	mockJWT.EXPECT().GenerateID().Times(1).
		Return(accessID)

	mockJWT.EXPECT().Assign(gomock.Any()).
		DoAndReturn(func(claims *jwtx.AAClaims) (string, error) {
			fmt.Println(claims.StdJWTClaims)
			assert.Equal(t, "http://localhost:8080", claims.StdJWTClaims.Audience)
			assert.Equal(t, accessID, claims.StdJWTClaims.Id)
			return newAccessToken, nil
		})

	mockOAuthRepo.EXPECT().SyncToken(ctx, gomock.Any()).
		DoAndReturn(func(c *gin.Context, token *model.Token) error {
			assert.Equal(t, newAccessToken, token.Access)
			assert.Equal(t, accessID, token.AccessID)
			return nil
		})

	svc := &service{
		oauthRepo: mockOAuthRepo,
		jwtx:      mockJWT,
	}

	resp, err := svc.Refresh(ctx, raw)
	assert.NoError(t, err)
	assert.Equal(t, &dto.RespCredential{
		Account:      jwtClaimSubject,
		Organization: jwtClaimIssuer,
		Environment:  "Horizon-Testnet",
		TokenPair: jwtx.TokenPair{
			Access:  newAccessToken,
			Refresh: raw,
		},
	}, resp)
}

func TestRefreshJWTParseFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockJWT := jwtx.NewMockJWT(ctrl)

	raw := "test-token"
	mockJWT.EXPECT().Parse(raw).Return(nil, errors.New("failed to parse JWT"))

	svc := &service{
		jwtx: mockJWT,
	}

	resp, err := svc.Refresh(ctx, raw)
	assert.Error(t, err)
	assert.Equal(t, "failed to parse JWT", err.Error())
	assert.Nil(t, resp)
}

func TestRefreshJWTParseClaimsFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockJWT := jwtx.NewMockJWT(ctrl)

	raw := "test-token"
	mockJWT.EXPECT().Parse(raw).Return(&jwt.StandardClaims{
		Id: "test-id",
	}, nil)

	svc := &service{
		jwtx: mockJWT,
	}

	resp, err := svc.Refresh(ctx, raw)
	assert.Error(t, err)
	assert.Equal(t, "failed to assert JWT claims as AAClaims", err.Error())
	assert.Nil(t, resp)
}

func TestRefreshFindTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockJWT := jwtx.NewMockJWT(ctrl)

	raw := "test-token"
	jwtClaimId := "test-id"

	mockJWT.EXPECT().Parse(raw).Return(&jwtx.AAClaims{
		StdJWTClaims: jwt.StandardClaims{
			Id: jwtClaimId,
		},
	}, nil)

	mockOAuthRepo.EXPECT().FindTokenByRefreshID(ctx, jwtClaimId).
		Return(nil, errors.New("failed to find token"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		jwtx:      mockJWT,
	}

	resp, err := svc.Refresh(ctx, raw)
	assert.Error(t, err)
	assert.Equal(t, "invalid refresh token", err.Error())
	assert.Nil(t, resp)
}

func TestRefreshAssignJWTError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockJWT := jwtx.NewMockJWT(ctrl)

	raw := "test-token"
	jwtClaimId := "test-id"
	accessID := "1"

	mockJWT.EXPECT().Parse(raw).Return(&jwtx.AAClaims{
		StdJWTClaims: jwt.StandardClaims{
			Id: jwtClaimId,
		},
	}, nil)

	mockOAuthRepo.EXPECT().FindTokenByRefreshID(ctx, jwtClaimId).
		Return(&model.Token{
			UserId: 1,
		}, nil)

	mockJWT.EXPECT().GenerateID().Times(1).
		Return(accessID)

	mockJWT.EXPECT().Assign(gomock.Any()).
		Return("", errors.New("failed to assign JWT"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		jwtx:      mockJWT,
	}

	resp, err := svc.Refresh(ctx, raw)
	assert.Error(t, err)
	assert.Equal(t, "failed to assign JWT", err.Error())
	assert.Nil(t, resp)
}

func TestRefreshSyncTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockJWT := jwtx.NewMockJWT(ctrl)

	raw := "test-token"
	jwtClaimId := "test-id"
	jwtClaimSubject := "account_name"
	jwtClaimIssuer := "org_name"
	newAccessToken := "new-access-token"
	accessID := "1"

	mockJWT.EXPECT().Parse(raw).Return(&jwtx.AAClaims{
		StdJWTClaims: jwt.StandardClaims{
			Id:      jwtClaimId,
			Subject: jwtClaimSubject,
			Issuer:  jwtClaimIssuer,
		},
	}, nil)

	mockOAuthRepo.EXPECT().FindTokenByRefreshID(ctx, jwtClaimId).
		Return(&model.Token{
			UserId: 1,
		}, nil)

	mockJWT.EXPECT().GenerateID().Times(1).
		Return(accessID)

	mockJWT.EXPECT().Assign(gomock.Any()).
		DoAndReturn(func(claims *jwtx.AAClaims) (string, error) {
			fmt.Println(claims.StdJWTClaims)
			assert.Equal(t, "http://localhost:8080", claims.StdJWTClaims.Audience)
			assert.Equal(t, accessID, claims.StdJWTClaims.Id)
			return newAccessToken, nil
		})

	mockOAuthRepo.EXPECT().SyncToken(ctx, gomock.Any()).
		Return(errors.New("failed to sync token"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		jwtx:      mockJWT,
	}

	resp, err := svc.Refresh(ctx, raw)
	assert.Error(t, err)
	assert.Equal(t, "failed to sync token", err.Error())
	assert.Nil(t, resp)
}

func TestLoginSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockDecrypter := decrypt.NewMockDecrypter(ctrl)
	mockJWT := jwtx.NewMockJWT(ctrl)

	rawPassword := []byte("raw_password")
	hashedPassword, _ := bcrypt.GenerateFromPassword(rawPassword, bcrypt.DefaultCost)
	accessID := "1"
	accountName := "account_name"
	orgName := "org_name"
	environment := "Horizon-Testnet"
	accessToken := "access_token"

	request := dto.ReqLogin{
		Organization: orgName,
		Account:      accountName,
		Password:     "password",
		Environment:  environment,
	}

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, req *dto.ReqOrgAcn) (*dto.RespUser, error) {
			assert.Equal(t, accountName, req.AcnName)
			assert.Equal(t, orgName, req.OrgName)
			return &dto.RespUser{
				Account:  accountName,
				Password: string(hashedPassword),
			}, nil
		})

	mockDecrypter.EXPECT().Decrypt(gomock.Any()).Times(1).
		DoAndReturn(func(data []byte) ([]byte, error) {
			assert.Equal(t, "password", string(data))
			return rawPassword, nil
		})

	mockJWT.EXPECT().GenerateID().Times(2).
		Return(accessID)

	mockJWT.EXPECT().Assign(gomock.Any()).Times(2).
		DoAndReturn(func(claims *jwtx.AAClaims) (string, error) {
			assert.Equal(t, accountName, claims.StdJWTClaims.Subject)
			assert.Equal(t, orgName, claims.StdJWTClaims.Issuer)
			assert.Equal(t, "http://localhost:8080", claims.StdJWTClaims.Audience)
			assert.Equal(t, accessID, claims.StdJWTClaims.Id)
			return accessToken, nil
		})

	mockOAuthRepo.EXPECT().SyncToken(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c *gin.Context, token *model.Token) error {
			assert.Equal(t, accessToken, token.Access)
			assert.Equal(t, accessID, token.AccessID)
			return nil
		})

	svc := &service{
		oauthRepo: mockOAuthRepo,
		decrypter: mockDecrypter,
		jwtx:      mockJWT,
	}

	resp, err := svc.Login(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, &dto.RespCredential{
		Account:      accountName,
		Organization: orgName,
		Environment:  environment,
		TokenPair: jwtx.TokenPair{
			Access:  accessToken,
			Refresh: accessToken,
		},
	}, resp)
}

func TestLoginFindUserByOrgAcnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)

	request := dto.ReqLogin{
		Organization: "org_name",
		Account:      "account_name",
	}

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(nil, errors.New("failed to find user by org and account"))
	svc := &service{
		oauthRepo: mockOAuthRepo,
	}

	resp, err := svc.Login(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "failed to find user by org and account", err.Error())
	assert.Nil(t, resp)
}

func TestLoginDecryptError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockDecrypter := decrypt.NewMockDecrypter(ctrl)

	request := dto.ReqLogin{
		Organization: "org_name",
		Account:      "account_name",
		Password:     "password",
	}

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			Account:  "account_name",
			Password: "hashed_password",
		}, nil)

	mockDecrypter.EXPECT().Decrypt(gomock.Any()).Times(1).
		Return(nil, errors.New("failed to decrypt password"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		decrypter: mockDecrypter,
	}

	resp, err := svc.Login(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "failed to decrypt password", err.Error())
	assert.Nil(t, resp)
}

func TestLoginComparePasswordError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockDecrypter := decrypt.NewMockDecrypter(ctrl)

	request := dto.ReqLogin{
		Organization: "org_name",
		Account:      "account_name",
		Password:     "password",
	}

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			Account:  "account_name",
			Password: "hashed_password",
		}, nil)

	mockDecrypter.EXPECT().Decrypt(gomock.Any()).Times(1).
		DoAndReturn(func(data []byte) ([]byte, error) {
			assert.Equal(t, "password", string(data))
			return []byte(""), nil
		})

	svc := &service{
		oauthRepo: mockOAuthRepo,
		decrypter: mockDecrypter,
	}

	resp, err := svc.Login(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "password not match", err.Error())
	assert.Nil(t, resp)
}

func TestLoginAssignJWTError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockDecrypter := decrypt.NewMockDecrypter(ctrl)
	mockJWT := jwtx.NewMockJWT(ctrl)

	rawPassword := []byte("raw_password")
	hashedPassword, _ := bcrypt.GenerateFromPassword(rawPassword, bcrypt.DefaultCost)
	accessID := "1"
	accountName := "account_name"
	orgName := "org_name"
	environment := "Horizon-Testnet"

	request := dto.ReqLogin{
		Organization: orgName,
		Account:      accountName,
		Password:     "password",
		Environment:  environment,
	}

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, req *dto.ReqOrgAcn) (*dto.RespUser, error) {
			assert.Equal(t, accountName, req.AcnName)
			assert.Equal(t, orgName, req.OrgName)
			return &dto.RespUser{
				Account:  accountName,
				Password: string(hashedPassword),
			}, nil
		})

	mockDecrypter.EXPECT().Decrypt(gomock.Any()).Times(1).
		DoAndReturn(func(data []byte) ([]byte, error) {
			assert.Equal(t, "password", string(data))
			return rawPassword, nil
		})

	mockJWT.EXPECT().GenerateID().Times(2).
		Return(accessID)

	mockJWT.EXPECT().Assign(gomock.Any()).Times(1).
		Return("", errors.New("failed to assign JWT"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		decrypter: mockDecrypter,
		jwtx:      mockJWT,
	}

	resp, err := svc.Login(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "failed to assign JWT", err.Error())
	assert.Nil(t, resp)
}

func TestLoginSyncTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockDecrypter := decrypt.NewMockDecrypter(ctrl)
	mockJWT := jwtx.NewMockJWT(ctrl)

	rawPassword := []byte("raw_password")
	hashedPassword, _ := bcrypt.GenerateFromPassword(rawPassword, bcrypt.DefaultCost)
	accessID := "1"
	accountName := "account_name"
	orgName := "org_name"
	environment := "Horizon-Testnet"
	accessToken := "access_token"

	request := dto.ReqLogin{
		Organization: orgName,
		Account:      accountName,
		Password:     "password",
		Environment:  environment,
	}

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, req *dto.ReqOrgAcn) (*dto.RespUser, error) {
			assert.Equal(t, accountName, req.AcnName)
			assert.Equal(t, orgName, req.OrgName)
			return &dto.RespUser{
				Account:  accountName,
				Password: string(hashedPassword),
			}, nil
		})

	mockDecrypter.EXPECT().Decrypt(gomock.Any()).Times(1).
		DoAndReturn(func(data []byte) ([]byte, error) {
			assert.Equal(t, "password", string(data))
			return rawPassword, nil
		})

	mockJWT.EXPECT().GenerateID().Times(2).
		Return(accessID)

	mockJWT.EXPECT().Assign(gomock.Any()).Times(2).
		DoAndReturn(func(claims *jwtx.AAClaims) (string, error) {
			assert.Equal(t, accountName, claims.StdJWTClaims.Subject)
			assert.Equal(t, orgName, claims.StdJWTClaims.Issuer)
			assert.Equal(t, "http://localhost:8080", claims.StdJWTClaims.Audience)
			assert.Equal(t, accessID, claims.StdJWTClaims.Id)
			return accessToken, nil
		})

	mockOAuthRepo.EXPECT().SyncToken(ctx, gomock.Any()).Times(1).
		Return(errors.New("failed to sync token"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		decrypter: mockDecrypter,
		jwtx:      mockJWT,
	}

	resp, err := svc.Login(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "failed to sync token", err.Error())
	assert.Nil(t, resp)
}

func TestSignupSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockDecrypter := decrypt.NewMockDecrypter(ctrl)
	mockJWT := jwtx.NewMockJWT(ctrl)
	mockResty := restyx.NewMockResty(ctrl)
	mockCsService := svcCS.NewMockService(ctrl)
	mockAmazon := amazonx.NewMockAmazon(ctrl)

	orgName := "org_name"
	accountName := "account_name"
	description := "description"
	request := dto.ReqSignup{
		Organization: orgName,
		Account:      accountName,
		Password:     "password",
		Description:  &description,
	}
	csToken := "cs_token"
	csRoleId := "cs_role_id"
	awsRoleArn := "arn:aws:iam::123456789012:role/AA-org_name-account_name-Role"
	awsRoleName := "AA-org_name-account_name-Role"
	awsSecretKey := "AA_org_name_account_name_SEC"
	awsSecretArn := "arn:aws:secretsmanager:us-east-1:123456789012:secret:AA_org_name_account_name_SEC-a1b2c3d4e5f6"

	mockOAuthRepo.EXPECT().FindOrgByName(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, name string) (*dto.RespOrg, error) {
			assert.Equal(t, orgName, name)
			return &dto.RespOrg{
				ID: 1,
			}, nil
		})

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, account string) (*dto.RespUser, error) {
			assert.Equal(t, accountName, account)
			return nil, errorx.NotFound("user/organization not found")
		})

	mockCsService.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockResty.EXPECT().AddCSRole(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, token string, orgName string, account string) (*dto.RespAddCsRole, error) {
			assert.Equal(t, csToken, token)
			assert.Equal(t, orgName, orgName)
			assert.Equal(t, account, account)
			return &dto.RespAddCsRole{
				RoleId: csRoleId,
			}, nil
		})

	mockAmazon.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
			assert.Equal(t, awsRoleName, *input.RoleName)
			return &iam.CreateRoleOutput{
				Role: &iamTypes.Role{
					Arn: aws.String(awsRoleArn),
				},
			}, nil
		})

	mockAmazon.EXPECT().PutRolePolicy(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
			assert.Equal(t, awsRoleName, *input.RoleName)
			assert.Equal(t, fmt.Sprintf("%s-policy", awsRoleName), *input.PolicyName)
			return nil, nil
		})

	mockAmazon.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *secretsmanager.CreateSecretInput) (*secretsmanager.CreateSecretOutput, error) {
			assert.Equal(t, awsSecretKey, *input.Name)
			assert.Equal(t, fmt.Sprintf(`{"cs_role": "%s"}`, csRoleId), *input.SecretString)
			return &secretsmanager.CreateSecretOutput{
				ARN: aws.String(awsSecretArn),
			}, nil
		})

	mockAmazon.EXPECT().PutResourcePolicy(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *secretsmanager.PutResourcePolicyInput) (*secretsmanager.PutResourcePolicyOutput, error) {
			assert.Equal(t, awsSecretKey, *input.SecretId)
			return nil, nil
		})

	mockDecrypter.EXPECT().Decrypt(gomock.Any()).Times(1).
		DoAndReturn(func(data []byte) ([]byte, error) {
			assert.Equal(t, "password", string(data))
			return []byte("hashed_password"), nil
		})

	mockOAuthRepo.EXPECT().CreateUser(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, user *model.User) error {
			assert.Equal(t, uint64(1), user.OrganizationId)
			assert.Equal(t, accountName, user.Account)
			assert.Equal(t, description, user.Description)
			return nil
		})

	svc := &service{
		oauthRepo: mockOAuthRepo,
		decrypter: mockDecrypter,
		jwtx:      mockJWT,
		resty:     mockResty,
		csService: mockCsService,
		amazon:    mockAmazon,
	}

	err := svc.Signup(ctx, request)
	assert.NoError(t, err)
}

func TestSignupFindOrgByNameError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)

	request := dto.ReqSignup{
		Organization: "org_name",
	}

	mockOAuthRepo.EXPECT().FindOrgByName(ctx, gomock.Any()).Times(1).
		Return(nil, errors.New("failed to find org by name"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
	}

	err := svc.Signup(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "failed to find org by name", err.Error())
}

func TestSignupFindUserByAcnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)

	request := dto.ReqSignup{
		Organization: "org_name",
		Account:      "account_name",
	}

	mockOAuthRepo.EXPECT().FindOrgByName(ctx, gomock.Any()).Times(1).
		Return(&dto.RespOrg{
			ID: 1,
		}, nil)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, gomock.Any()).Times(1).
		Return(nil, errors.New("failed to find user by account name"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
	}

	err := svc.Signup(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "failed to find user by account name", err.Error())
}

func TestSignupUserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)

	request := dto.ReqSignup{
		Organization: "org_name",
		Account:      "account_name",
	}

	mockOAuthRepo.EXPECT().FindOrgByName(ctx, gomock.Any()).Times(1).
		Return(&dto.RespOrg{
			ID: 1,
		}, nil)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			Account: "account_name",
		}, nil)

	svc := &service{
		oauthRepo: mockOAuthRepo,
	}

	err := svc.Signup(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "user already exists", err.Error())
}

func TestSignupGetCSTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCsService := svcCS.NewMockService(ctrl)

	accountName := "account_name"
	request := dto.ReqSignup{
		Organization: "org_name",
		Account:      accountName,
	}

	mockOAuthRepo.EXPECT().FindOrgByName(ctx, gomock.Any()).Times(1).
		Return(&dto.RespOrg{
			ID: 1,
		}, nil)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, account string) (*dto.RespUser, error) {
			assert.Equal(t, accountName, account)
			return nil, errorx.NotFound("user/organization not found")
		})

	mockCsService.EXPECT().CubeSignerToken(ctx).Times(1).
		Return("", errors.New("failed to get cube signer token"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csService: mockCsService,
	}

	err := svc.Signup(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "failed to get cube signer token", err.Error())
}

func TestSignupAddCSRoleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCsService := svcCS.NewMockService(ctrl)
	mockResty := restyx.NewMockResty(ctrl)

	orgName := "org_name"
	accountName := "account_name"
	request := dto.ReqSignup{
		Organization: orgName,
		Account:      accountName,
	}
	csToken := "cs_token"

	mockOAuthRepo.EXPECT().FindOrgByName(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, name string) (*dto.RespOrg, error) {
			assert.Equal(t, orgName, name)
			return &dto.RespOrg{
				ID: 1,
			}, nil
		})

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, account string) (*dto.RespUser, error) {
			assert.Equal(t, accountName, account)
			return nil, errorx.NotFound("user/organization not found")
		})

	mockCsService.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockResty.EXPECT().AddCSRole(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
		Return(nil, errors.New("failed to add cs role"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csService: mockCsService,
		resty:     mockResty,
	}

	err := svc.Signup(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "failed to add cs role", err.Error())
}

func TestSignupAddAwsRoleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockResty := restyx.NewMockResty(ctrl)
	mockCsService := svcCS.NewMockService(ctrl)
	mockAmazon := amazonx.NewMockAmazon(ctrl)

	orgName := "org_name"
	accountName := "account_name"
	description := "description"
	request := dto.ReqSignup{
		Organization: orgName,
		Account:      accountName,
		Password:     "password",
		Description:  &description,
	}
	csToken := "cs_token"
	csRoleId := "cs_role_id"

	mockOAuthRepo.EXPECT().FindOrgByName(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, name string) (*dto.RespOrg, error) {
			assert.Equal(t, orgName, name)
			return &dto.RespOrg{
				ID: 1,
			}, nil
		})

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, account string) (*dto.RespUser, error) {
			assert.Equal(t, accountName, account)
			return nil, errorx.NotFound("user/organization not found")
		})

	mockCsService.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockResty.EXPECT().AddCSRole(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, token string, orgName string, account string) (*dto.RespAddCsRole, error) {
			assert.Equal(t, csToken, token)
			assert.Equal(t, orgName, orgName)
			assert.Equal(t, account, account)
			return &dto.RespAddCsRole{
				RoleId: csRoleId,
			}, nil
		})

	mockAmazon.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Times(1).
		Return(nil, errors.New("failed to create aws role"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		resty:     mockResty,
		csService: mockCsService,
		amazon:    mockAmazon,
	}

	err := svc.Signup(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "create aws role occurred error: failed to create aws role", err.Error())
}

func TestSignupPutRolePolicyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockResty := restyx.NewMockResty(ctrl)
	mockCsService := svcCS.NewMockService(ctrl)
	mockAmazon := amazonx.NewMockAmazon(ctrl)

	orgName := "org_name"
	accountName := "account_name"
	description := "description"
	request := dto.ReqSignup{
		Organization: orgName,
		Account:      accountName,
		Password:     "password",
		Description:  &description,
	}
	csToken := "cs_token"
	csRoleId := "cs_role_id"
	awsRoleName := "AA-org_name-account_name-Role"
	awsRoleArn := "arn:aws:iam::123456789012:role/AA-org_name-account_name-Role"

	mockOAuthRepo.EXPECT().FindOrgByName(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, name string) (*dto.RespOrg, error) {
			assert.Equal(t, orgName, name)
			return &dto.RespOrg{
				ID: 1,
			}, nil
		})

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, account string) (*dto.RespUser, error) {
			assert.Equal(t, accountName, account)
			return nil, errorx.NotFound("user/organization not found")
		})

	mockCsService.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockResty.EXPECT().AddCSRole(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, token string, orgName string, account string) (*dto.RespAddCsRole, error) {
			assert.Equal(t, csToken, token)
			assert.Equal(t, orgName, orgName)
			assert.Equal(t, account, account)
			return &dto.RespAddCsRole{
				RoleId: csRoleId,
			}, nil
		})

	mockAmazon.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
			assert.Equal(t, awsRoleName, *input.RoleName)
			return &iam.CreateRoleOutput{
				Role: &iamTypes.Role{
					Arn: aws.String(awsRoleArn),
				},
			}, nil
		})

	mockAmazon.EXPECT().PutRolePolicy(gomock.Any(), gomock.Any()).Times(1).
		Return(nil, errors.New("failed to put role policy"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		resty:     mockResty,
		csService: mockCsService,
		amazon:    mockAmazon,
	}

	err := svc.Signup(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "put role policy for aws role occurred error: failed to put role policy", err.Error())
}

func TestSignupCreateSecretError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockDecrypter := decrypt.NewMockDecrypter(ctrl)
	mockJWT := jwtx.NewMockJWT(ctrl)
	mockResty := restyx.NewMockResty(ctrl)
	mockCsService := svcCS.NewMockService(ctrl)
	mockAmazon := amazonx.NewMockAmazon(ctrl)

	orgName := "org_name"
	accountName := "account_name"
	description := "description"
	request := dto.ReqSignup{
		Organization: orgName,
		Account:      accountName,
		Password:     "password",
		Description:  &description,
	}
	csToken := "cs_token"
	csRoleId := "cs_role_id"
	awsRoleArn := "arn:aws:iam::123456789012:role/AA-org_name-account_name-Role"
	awsRoleName := "AA-org_name-account_name-Role"

	mockOAuthRepo.EXPECT().FindOrgByName(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, name string) (*dto.RespOrg, error) {
			assert.Equal(t, orgName, name)
			return &dto.RespOrg{
				ID: 1,
			}, nil
		})

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, account string) (*dto.RespUser, error) {
			assert.Equal(t, accountName, account)
			return nil, errorx.NotFound("user/organization not found")
		})

	mockCsService.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockResty.EXPECT().AddCSRole(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, csToken string, orgName string, account string) (*dto.RespAddCsRole, error) {
			assert.Equal(t, csToken, csToken)
			assert.Equal(t, orgName, orgName)
			assert.Equal(t, account, account)
			return &dto.RespAddCsRole{
				RoleId: csRoleId,
			}, nil
		})

	mockAmazon.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
			assert.Equal(t, awsRoleName, *input.RoleName)
			return &iam.CreateRoleOutput{
				Role: &iamTypes.Role{
					Arn: aws.String(awsRoleArn),
				},
			}, nil
		})

	mockAmazon.EXPECT().PutRolePolicy(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
			assert.Equal(t, awsRoleName, *input.RoleName)
			assert.Equal(t, fmt.Sprintf("%s-policy", awsRoleName), *input.PolicyName)
			return nil, nil
		})

	mockAmazon.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).Times(1).
		Return(nil, errors.New("failed to create secret"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		decrypter: mockDecrypter,
		jwtx:      mockJWT,
		resty:     mockResty,
		csService: mockCsService,
		amazon:    mockAmazon,
	}

	err := svc.Signup(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "create aws secret key occurred error: failed to create secret", err.Error())
}

func TestSignupPutResourcePolicyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockDecrypter := decrypt.NewMockDecrypter(ctrl)
	mockJWT := jwtx.NewMockJWT(ctrl)
	mockResty := restyx.NewMockResty(ctrl)
	mockCsService := svcCS.NewMockService(ctrl)
	mockAmazon := amazonx.NewMockAmazon(ctrl)

	orgName := "org_name"
	accountName := "account_name"
	description := "description"
	request := dto.ReqSignup{
		Organization: orgName,
		Account:      accountName,
		Password:     "password",
		Description:  &description,
	}
	csToken := "cs_token"
	csRoleId := "cs_role_id"
	awsRoleArn := "arn:aws:iam::123456789012:role/AA-org_name-account_name-Role"
	awsRoleName := "AA-org_name-account_name-Role"
	awsSecretArn := "arn:aws:secretsmanager:us-east-1:123456789012:secret:AA_org_name_account_name_SEC-a1b2c3d4e5f6"

	mockOAuthRepo.EXPECT().FindOrgByName(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, name string) (*dto.RespOrg, error) {
			assert.Equal(t, orgName, name)
			return &dto.RespOrg{
				ID: 1,
			}, nil
		})

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, account string) (*dto.RespUser, error) {
			assert.Equal(t, accountName, account)
			return nil, errorx.NotFound("user/organization not found")
		})

	mockCsService.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockResty.EXPECT().AddCSRole(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, token string, orgName string, account string) (*dto.RespAddCsRole, error) {
			assert.Equal(t, csToken, token)
			assert.Equal(t, orgName, orgName)
			assert.Equal(t, account, account)
			return &dto.RespAddCsRole{
				RoleId: csRoleId,
			}, nil
		})

	mockAmazon.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
			assert.Equal(t, awsRoleName, *input.RoleName)
			return &iam.CreateRoleOutput{
				Role: &iamTypes.Role{
					Arn: aws.String(awsRoleArn),
				},
			}, nil
		})

	mockAmazon.EXPECT().PutRolePolicy(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
			assert.Equal(t, awsRoleName, *input.RoleName)
			assert.Equal(t, fmt.Sprintf("%s-policy", awsRoleName), *input.PolicyName)
			return nil, nil
		})

	mockAmazon.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *secretsmanager.CreateSecretInput) (*secretsmanager.CreateSecretOutput, error) {
			return &secretsmanager.CreateSecretOutput{
				ARN: aws.String(awsSecretArn),
			}, nil
		})

	mockAmazon.EXPECT().PutResourcePolicy(gomock.Any(), gomock.Any()).Times(1).
		Return(nil, errors.New("failed to put resource policy"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		decrypter: mockDecrypter,
		jwtx:      mockJWT,
		resty:     mockResty,
		csService: mockCsService,
		amazon:    mockAmazon,
	}

	err := svc.Signup(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "put resource policy for aws secret key occurred error: failed to put resource policy", err.Error())
}

func TestSignupDecryptPwdError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockDecrypter := decrypt.NewMockDecrypter(ctrl)
	mockJWT := jwtx.NewMockJWT(ctrl)
	mockResty := restyx.NewMockResty(ctrl)
	mockCsService := svcCS.NewMockService(ctrl)
	mockAmazon := amazonx.NewMockAmazon(ctrl)

	orgName := "org_name"
	accountName := "account_name"
	description := "description"
	request := dto.ReqSignup{
		Organization: orgName,
		Account:      accountName,
		Password:     "password",
		Description:  &description,
	}
	csToken := "cs_token"
	csRoleId := "cs_role_id"
	awsRoleArn := "arn:aws:iam::123456789012:role/AA-org_name-account_name-Role"
	awsRoleName := "AA-org_name-account_name-Role"
	awsSecretKey := "AA_org_name_account_name_SEC"
	awsSecretArn := "arn:aws:secretsmanager:us-east-1:123456789012:secret:AA_org_name_account_name_SEC-a1b2c3d4e5f6"

	mockOAuthRepo.EXPECT().FindOrgByName(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, name string) (*dto.RespOrg, error) {
			assert.Equal(t, orgName, name)
			return &dto.RespOrg{
				ID: 1,
			}, nil
		})

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, account string) (*dto.RespUser, error) {
			assert.Equal(t, accountName, account)
			return nil, errorx.NotFound("user/organization not found")
		})

	mockCsService.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockResty.EXPECT().AddCSRole(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, csToken string, orgName string, account string) (*dto.RespAddCsRole, error) {
			assert.Equal(t, csToken, csToken)
			assert.Equal(t, orgName, orgName)
			assert.Equal(t, account, account)
			return &dto.RespAddCsRole{
				RoleId: csRoleId,
			}, nil
		})

	mockAmazon.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
			assert.Equal(t, awsRoleName, *input.RoleName)
			return &iam.CreateRoleOutput{
				Role: &iamTypes.Role{
					Arn: aws.String(awsRoleArn),
				},
			}, nil
		})

	mockAmazon.EXPECT().PutRolePolicy(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
			assert.Equal(t, awsRoleName, *input.RoleName)
			assert.Equal(t, fmt.Sprintf("%s-policy", awsRoleName), *input.PolicyName)
			return nil, nil
		})

	mockAmazon.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *secretsmanager.CreateSecretInput) (*secretsmanager.CreateSecretOutput, error) {
			assert.Equal(t, awsSecretKey, *input.Name)
			assert.Equal(t, fmt.Sprintf(`{"cs_role": "%s"}`, csRoleId), *input.SecretString)
			return &secretsmanager.CreateSecretOutput{
				ARN: aws.String(awsSecretArn),
			}, nil
		})

	mockAmazon.EXPECT().PutResourcePolicy(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *secretsmanager.PutResourcePolicyInput) (*secretsmanager.PutResourcePolicyOutput, error) {
			assert.Equal(t, awsSecretKey, *input.SecretId)
			return nil, nil
		})

	mockDecrypter.EXPECT().Decrypt(gomock.Any()).Times(1).
		Return(nil, errors.New("failed to decrypt"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		decrypter: mockDecrypter,
		jwtx:      mockJWT,
		resty:     mockResty,
		csService: mockCsService,
		amazon:    mockAmazon,
	}

	err := svc.Signup(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "failed to decrypt", err.Error())
}

func TestSignupCreateUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockDecrypter := decrypt.NewMockDecrypter(ctrl)
	mockJWT := jwtx.NewMockJWT(ctrl)
	mockResty := restyx.NewMockResty(ctrl)
	mockCsService := svcCS.NewMockService(ctrl)
	mockAmazon := amazonx.NewMockAmazon(ctrl)

	orgName := "org_name"
	accountName := "account_name"
	description := "description"
	request := dto.ReqSignup{
		Organization: orgName,
		Account:      accountName,
		Password:     "password",
		Description:  &description,
	}
	csToken := "cs_token"
	csRoleId := "cs_role_id"
	awsRoleArn := "arn:aws:iam::123456789012:role/AA-org_name-account_name-Role"
	awsRoleName := "AA-org_name-account_name-Role"
	awsSecretKey := "AA_org_name_account_name_SEC"
	awsSecretArn := "arn:aws:secretsmanager:us-east-1:123456789012:secret:AA_org_name_account_name_SEC-a1b2c3d4e5f6"

	mockOAuthRepo.EXPECT().FindOrgByName(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, name string) (*dto.RespOrg, error) {
			assert.Equal(t, orgName, name)
			return &dto.RespOrg{
				ID: 1,
			}, nil
		})

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, account string) (*dto.RespUser, error) {
			assert.Equal(t, accountName, account)
			return nil, errorx.NotFound("user/organization not found")
		})

	mockCsService.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockResty.EXPECT().AddCSRole(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, csToken string, orgName string, account string) (*dto.RespAddCsRole, error) {
			assert.Equal(t, csToken, csToken)
			assert.Equal(t, orgName, orgName)
			assert.Equal(t, account, account)
			return &dto.RespAddCsRole{
				RoleId: csRoleId,
			}, nil
		})

	mockAmazon.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
			assert.Equal(t, awsRoleName, *input.RoleName)
			return &iam.CreateRoleOutput{
				Role: &iamTypes.Role{
					Arn: aws.String(awsRoleArn),
				},
			}, nil
		})

	mockAmazon.EXPECT().PutRolePolicy(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
			assert.Equal(t, awsRoleName, *input.RoleName)
			assert.Equal(t, fmt.Sprintf("%s-policy", awsRoleName), *input.PolicyName)
			return nil, nil
		})

	mockAmazon.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *secretsmanager.CreateSecretInput) (*secretsmanager.CreateSecretOutput, error) {
			assert.Equal(t, awsSecretKey, *input.Name)
			assert.Equal(t, fmt.Sprintf(`{"cs_role": "%s"}`, csRoleId), *input.SecretString)
			return &secretsmanager.CreateSecretOutput{
				ARN: aws.String(awsSecretArn),
			}, nil
		})

	mockAmazon.EXPECT().PutResourcePolicy(gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, input *secretsmanager.PutResourcePolicyInput) (*secretsmanager.PutResourcePolicyOutput, error) {
			assert.Equal(t, awsSecretKey, *input.SecretId)
			return nil, nil
		})

	mockDecrypter.EXPECT().Decrypt(gomock.Any()).Times(1).
		DoAndReturn(func(data []byte) ([]byte, error) {
			assert.Equal(t, "password", string(data))
			return []byte("hashed_password"), nil
		})

	mockOAuthRepo.EXPECT().CreateUser(ctx, gomock.Any()).Times(1).
		Return(errors.New("failed to create user"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		decrypter: mockDecrypter,
		jwtx:      mockJWT,
		resty:     mockResty,
		csService: mockCsService,
		amazon:    mockAmazon,
	}

	err := svc.Signup(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "failed to create user", err.Error())
}
