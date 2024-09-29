package wallet

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/constant"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/repo"
	svcCS "github.com/57blocks/auto-action/server/internal/service/cs"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"
	"github.com/57blocks/auto-action/server/internal/third-party/restyx"
	"github.com/57blocks/auto-action/server/internal/third-party/stellarx"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Before test, setup log and config
func TestMain(m *testing.M) {
	os.Setenv("WALLET_MAX", "1")
	os.Setenv("BOUND_NAME", "Horizon-Testnet")
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

func TestCreateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	csToken := "cs-token"
	testKey := "test-key"
	testRole := "test-role"
	testCSKey := "Key#Stellar_test-key"

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockResty := restyx.NewMockResty(ctrl)
	mockCS := svcCS.NewMockService(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(c context.Context, req *dto.ReqOrgAcn) (*dto.RespUser, error) {
			assert.Equal(t, testOrg, req.OrgName)
			assert.Equal(t, testAccount, req.AcnName)
			return &dto.RespUser{
				ID: 1,
			}, nil
		})

	mockCSRepo.EXPECT().FindCSKeysByAccount(ctx, uint64(1)).Times(1).
		Return([]*model.CubeSignerKey{}, nil)

	mockCS.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockCS.EXPECT().GetSecRole(ctx, "AA_test-org_test-account_SEC").Times(1).
		Return(testRole, nil)

	mockResty.EXPECT().AddCSKey(ctx, csToken).Times(1).
		Return(testCSKey, nil)

	mockResty.EXPECT().AddCSKeyToRole(ctx, csToken, testCSKey, testRole).Times(1).
		Return(nil)

	mockCSRepo.EXPECT().SyncCSKey(ctx, gomock.Any()).Times(1).
		Return(nil)

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		resty:     mockResty,
		csService: mockCS,
	}

	wallet, err := svc.Create(ctx)
	assert.NoError(t, err)
	assert.Equal(t, testKey, wallet.Address)
}

func TestCreateFindUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)

	mockOAuthRepo := repo.NewMockOAuth(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(nil, errors.New("user not found"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
	}

	wallet, err := svc.Create(ctx)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, wallet)
}

func TestCreateMaxWalletError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKeysByAccount(ctx, uint64(1)).Times(1).
		Return([]*model.CubeSignerKey{
			{
				Key: "test-key",
			},
		}, nil)

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
	}

	wallet, err := svc.Create(ctx)
	assert.Error(t, err)
	assert.Equal(t, "the number of wallet address is limited to 1", err.Error())
	assert.Nil(t, wallet)
}

func TestCreateCubeSignerTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockCS := svcCS.NewMockService(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKeysByAccount(ctx, uint64(1)).Times(1).
		Return([]*model.CubeSignerKey{}, nil)

	mockCS.EXPECT().CubeSignerToken(ctx).Times(1).
		Return("", errors.New("cube signer token error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		csService: mockCS,
	}

	wallet, err := svc.Create(ctx)
	assert.Error(t, err)
	assert.Equal(t, "cube signer token error", err.Error())
	assert.Nil(t, wallet)
}

func TestCreateGetSecRoleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	csToken := "cs-token"

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockCS := svcCS.NewMockService(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKeysByAccount(ctx, uint64(1)).Times(1).
		Return([]*model.CubeSignerKey{}, nil)

	mockCS.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockCS.EXPECT().GetSecRole(ctx, "AA_test-org_test-account_SEC").Times(1).
		Return("", errors.New("get sec role error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		csService: mockCS,
	}

	wallet, err := svc.Create(ctx)
	assert.Error(t, err)
	assert.Equal(t, "get sec role error", err.Error())
	assert.Nil(t, wallet)
}

func TestCreateAddCSKeyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	csToken := "cs-token"

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockCS := svcCS.NewMockService(ctrl)
	mockResty := restyx.NewMockResty(ctrl)
	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKeysByAccount(ctx, uint64(1)).Times(1).
		Return([]*model.CubeSignerKey{}, nil)

	mockCS.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockCS.EXPECT().GetSecRole(ctx, "AA_test-org_test-account_SEC").Times(1).
		Return("test-role", nil)

	mockResty.EXPECT().AddCSKey(ctx, csToken).Times(1).
		Return("", errors.New("add cs key error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		csService: mockCS,
		resty:     mockResty,
	}

	wallet, err := svc.Create(ctx)
	assert.Error(t, err)
	assert.Equal(t, "add cs key error", err.Error())
	assert.Nil(t, wallet)
}

func TestCreateAddCSKeyToRoleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	csToken := "cs-token"
	testKey := "test-key"
	testRole := "test-role"

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockCS := svcCS.NewMockService(ctrl)
	mockResty := restyx.NewMockResty(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKeysByAccount(ctx, uint64(1)).Times(1).
		Return([]*model.CubeSignerKey{}, nil)

	mockCS.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockCS.EXPECT().GetSecRole(ctx, "AA_test-org_test-account_SEC").Times(1).
		Return(testRole, nil)

	mockResty.EXPECT().AddCSKey(ctx, csToken).Times(1).
		Return(testKey, nil)

	mockResty.EXPECT().AddCSKeyToRole(ctx, csToken, testKey, testRole).Times(1).
		Return(errors.New("add cs key to role error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		csService: mockCS,
		resty:     mockResty,
	}

	wallet, err := svc.Create(ctx)
	assert.Error(t, err)
	assert.Equal(t, "add cs key to role error", err.Error())
	assert.Nil(t, wallet)
}

func TestCreateSyncCSKeyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	csToken := "cs-token"
	testKey := "test-key"
	testRole := "test-role"

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockCS := svcCS.NewMockService(ctrl)
	mockResty := restyx.NewMockResty(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKeysByAccount(ctx, uint64(1)).Times(1).
		Return([]*model.CubeSignerKey{}, nil)

	mockCS.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockCS.EXPECT().GetSecRole(ctx, "AA_test-org_test-account_SEC").Times(1).
		Return(testRole, nil)

	mockResty.EXPECT().AddCSKey(ctx, csToken).Times(1).
		Return(testKey, nil)

	mockResty.EXPECT().AddCSKeyToRole(ctx, csToken, testKey, testRole).Times(1).
		Return(nil)

	mockCSRepo.EXPECT().SyncCSKey(ctx, gomock.Any()).Times(1).
		Return(errors.New("sync cs key error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		csService: mockCS,
		resty:     mockResty,
	}

	wallet, err := svc.Create(ctx)
	assert.Error(t, err)
	assert.Equal(t, "sync cs key error", err.Error())
	assert.Nil(t, wallet)
}

func TestRemoveSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	testKeyId := "Key#Stellar_test-key"
	csToken := "cs-token"
	testRole := "test-role"

	request := &dto.ReqRemoveWallet{
		Address: "test-key",
	}

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockCS := svcCS.NewMockService(ctrl)
	mockResty := restyx.NewMockResty(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKey(ctx, testKeyId, uint64(1)).Times(1).
		Return(&model.CubeSignerKey{
			Key: testKeyId,
		}, nil)

	mockCS.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockCS.EXPECT().GetSecRole(ctx, "AA_test-org_test-account_SEC").Times(1).
		Return(testRole, nil)

	mockResty.EXPECT().DeleteCSKeyFromRole(ctx, csToken, testKeyId, testRole).Times(1).
		Return(nil)

	mockResty.EXPECT().DeleteCSKey(ctx, csToken, testKeyId).Times(1).
		Return(nil)

	mockCSRepo.EXPECT().DeleteCSKey(ctx, testKeyId, uint64(1)).Times(1).
		Return(nil)

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		csService: mockCS,
		resty:     mockResty,
	}

	err := svc.Remove(ctx, request)
	assert.NoError(t, err)
}

func TestRemoveFindUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	request := &dto.ReqRemoveWallet{
		Address: "test-key",
	}

	mockOAuthRepo := repo.NewMockOAuth(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(nil, errors.New("user not found"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
	}

	err := svc.Remove(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestRemoveFindCSKeyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	request := &dto.ReqRemoveWallet{
		Address: "test-key",
	}

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKey(ctx, gomock.Any(), uint64(1)).Times(1).
		Return(nil, errors.New("cube signer key not found"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
	}

	err := svc.Remove(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "no existed wallet address found: test-key", err.Error())
}

func TestRemoveCubeSignerTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	request := &dto.ReqRemoveWallet{
		Address: "test-key",
	}

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockCS := svcCS.NewMockService(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKey(ctx, gomock.Any(), uint64(1)).Times(1).
		Return(&model.CubeSignerKey{
			Key: "test-key",
		}, nil)

	mockCS.EXPECT().CubeSignerToken(ctx).Times(1).
		Return("", errors.New("cube signer token error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		csService: mockCS,
	}

	err := svc.Remove(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "cube signer token error", err.Error())
}

func TestRemoveGetSecRoleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	csToken := "cs-token"
	request := &dto.ReqRemoveWallet{
		Address: "test-key",
	}

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockCS := svcCS.NewMockService(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKey(ctx, gomock.Any(), uint64(1)).Times(1).
		Return(&model.CubeSignerKey{
			Key: "test-key",
		}, nil)

	mockCS.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockCS.EXPECT().GetSecRole(ctx, "AA_test-org_test-account_SEC").Times(1).
		Return("", errors.New("get sec role error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		csService: mockCS,
	}

	err := svc.Remove(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "get sec role error", err.Error())
}

func TestRemoveDeleteCSKeyFromRoleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	csToken := "cs-token"
	testKeyId := "Key#Stellar_test-key"
	testRole := "test-role"
	request := &dto.ReqRemoveWallet{
		Address: "test-key",
	}

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockCS := svcCS.NewMockService(ctrl)
	mockResty := restyx.NewMockResty(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKey(ctx, gomock.Any(), uint64(1)).Times(1).
		Return(&model.CubeSignerKey{
			Key: testKeyId,
		}, nil)

	mockCS.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockCS.EXPECT().GetSecRole(ctx, "AA_test-org_test-account_SEC").Times(1).
		Return(testRole, nil)

	mockResty.EXPECT().DeleteCSKeyFromRole(ctx, csToken, testKeyId, testRole).Times(1).
		Return(errors.New("delete cs key from role error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		csService: mockCS,
		resty:     mockResty,
	}

	err := svc.Remove(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "delete cs key from role error", err.Error())
}

func TestRemoveDeleteCSKeyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	csToken := "cs-token"
	testKeyId := "Key#Stellar_test-key"
	testRole := "test-role"
	request := &dto.ReqRemoveWallet{
		Address: "test-key",
	}

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockCS := svcCS.NewMockService(ctrl)
	mockResty := restyx.NewMockResty(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKey(ctx, gomock.Any(), uint64(1)).Times(1).
		Return(&model.CubeSignerKey{
			Key: testKeyId,
		}, nil)

	mockCS.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockCS.EXPECT().GetSecRole(ctx, "AA_test-org_test-account_SEC").Times(1).
		Return(testRole, nil)

	mockResty.EXPECT().DeleteCSKeyFromRole(ctx, csToken, testKeyId, testRole).Times(1).
		Return(nil)

	mockResty.EXPECT().DeleteCSKey(ctx, csToken, testKeyId).Times(1).
		Return(errors.New("delete cs key error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		csService: mockCS,
		resty:     mockResty,
	}

	err := svc.Remove(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "delete cs key error", err.Error())
}

func TestRemoveDeleteRepoCSKeyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	csToken := "cs-token"
	testKeyId := "Key#Stellar_test-key"
	testRole := "test-role"
	request := &dto.ReqRemoveWallet{
		Address: "test-key",
	}

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockCS := svcCS.NewMockService(ctrl)
	mockResty := restyx.NewMockResty(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKey(ctx, gomock.Any(), uint64(1)).Times(1).
		Return(&model.CubeSignerKey{
			Key: testKeyId,
		}, nil)

	mockCS.EXPECT().CubeSignerToken(ctx).Times(1).
		Return(csToken, nil)

	mockCS.EXPECT().GetSecRole(ctx, "AA_test-org_test-account_SEC").Times(1).
		Return(testRole, nil)

	mockResty.EXPECT().DeleteCSKeyFromRole(ctx, csToken, testKeyId, testRole).Times(1).
		Return(nil)

	mockResty.EXPECT().DeleteCSKey(ctx, csToken, testKeyId).Times(1).
		Return(nil)

	mockCSRepo.EXPECT().DeleteCSKey(ctx, testKeyId, uint64(1)).Times(1).
		Return(errors.New("delete repo cs key error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		csService: mockCS,
		resty:     mockResty,
	}

	err := svc.Remove(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "delete repo cs key error", err.Error())
}

func TestListSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	testKeyId := "Key#Stellar_test-key"

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKeysByAccount(ctx, uint64(1)).Times(1).
		Return([]*model.CubeSignerKey{
			{
				Key: testKeyId,
			},
		}, nil)

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
	}

	wallets, err := svc.List(ctx)
	assert.NoError(t, err)
	assert.Equal(t, &dto.RespListWallets{
		Data: []dto.RespListWallet{
			{
				Address: "test-key",
			},
		},
	}, wallets)
}

func TestListFindUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)

	mockOAuthRepo := repo.NewMockOAuth(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(nil, errors.New("user not found"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
	}

	wallets, err := svc.List(ctx)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, wallets)
}

func TestListFindCSKeysError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKeysByAccount(ctx, uint64(1)).Times(1).
		Return(nil, errors.New("find cs keys error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
	}

	wallets, err := svc.List(ctx)
	assert.Error(t, err)
	assert.Equal(t, "find cs keys error", err.Error())
	assert.Nil(t, wallets)
}

func TestVerifySuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	testKeyId := "Key#Stellar_test-key"
	testAddress := "test-key"
	request := &dto.ReqVerifyWallet{
		Address: testAddress,
	}

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockStellar := stellarx.NewMockStellar(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKey(ctx, testKeyId, uint64(1)).Times(1).
		Return(&model.CubeSignerKey{
			Key: testKeyId,
		}, nil)

	mockStellar.EXPECT().AccountDetail(ctx, horizonclient.AccountRequest{AccountID: testAddress}).Times(1).
		Return(horizon.Account{}, nil)

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		stellar:   mockStellar,
	}

	wallet, err := svc.Verify(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, &dto.RespVerifyWallet{
		Address: testAddress,
		IsValid: true,
	}, wallet)
}

func TestVerifyFindUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	testAddress := "test-key"
	request := &dto.ReqVerifyWallet{
		Address: testAddress,
	}

	mockOAuthRepo := repo.NewMockOAuth(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(nil, errors.New("user not found"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
	}

	wallet, err := svc.Verify(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, wallet)
}

func TestVerifyFindCSKeyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	testKeyId := "Key#Stellar_test-key"
	testAddress := "test-key"
	request := &dto.ReqVerifyWallet{
		Address: testAddress,
	}

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKey(ctx, testKeyId, uint64(1)).Times(1).
		Return(nil, errors.New("find cs key error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
	}

	wallet, err := svc.Verify(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "find cs key error", err.Error())
	assert.Nil(t, wallet)
}

func TestVerifyAccountDetailError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	testOrg := "test-org"
	testAccount := "test-account"
	ctx.Set(constant.ClaimIss.Str(), testOrg)
	ctx.Set(constant.ClaimSub.Str(), testAccount)
	testKeyId := "Key#Stellar_test-key"
	testAddress := "test-key"
	request := &dto.ReqVerifyWallet{
		Address: testAddress,
	}

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)
	mockStellar := stellarx.NewMockStellar(ctrl)

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, gomock.Any()).Times(1).
		Return(&dto.RespUser{
			ID: 1,
		}, nil)

	mockCSRepo.EXPECT().FindCSKey(ctx, testKeyId, uint64(1)).Times(1).
		Return(&model.CubeSignerKey{
			Key: testKeyId,
		}, nil)

	mockStellar.EXPECT().AccountDetail(ctx, horizonclient.AccountRequest{AccountID: testAddress}).Times(1).
		Return(horizon.Account{}, errors.New("account detail error"))

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
		stellar:   mockStellar,
	}

	wallet, err := svc.Verify(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, &dto.RespVerifyWallet{
		Address: testAddress,
		IsValid: false,
	}, wallet)
}
