package cs

import (
	"context"
	"fmt"
	"testing"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/repo"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestToSignSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)

	ctx := context.TODO()

	userID := uint64(1)
	from := "from"
	account := "v3n"
	organization := "epic"

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, &dto.ReqOrgAcn{
		OrgName: "epic",
		AcnName: "v3n",
	}).Times(1).
		Return(
			&dto.RespUser{
				ID:             userID,
				Account:        account,
				Password:       "",
				Description:    "",
				CubeSignerUser: "",
				OrganizationId: 1,
				Organization: &dto.RespOrg{
					ID:          1,
					Name:        organization,
					Description: "",
				},
			}, nil,
		)

	mockCSRepo.EXPECT().ToSign(ctx, userID, fmt.Sprintf("Key#Stellar_%s", from)).Times(1).
		Return(
			&dto.RespCSKey{
				Account: dto.RespUser{
					Account:        account,
					CubeSignerUser: "User#_",
				},
				Organization: "Org#_",
				Key:          "Key#Stellar_ABCDEFG",
				Scopes:       []string{"sign:blob"},
			}, nil,
		)

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
	}
	toSign, err := svc.ToSign(ctx, &dto.ReqToSign{Account: account, Organization: organization, From: from})
	assert.NoError(t, err)
	assert.Equal(t, "Key#Stellar_ABCDEFG", toSign.Key)
}

func TestToSignWithKeyNotFoundError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)

	ctx := context.TODO()

	userID := uint64(1)
	from := "from"
	account := "v3n"
	organization := "epic"

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, &dto.ReqOrgAcn{
		OrgName: "epic",
		AcnName: "v3n",
	}).Times(1).
		Return(
			&dto.RespUser{
				ID:             userID,
				Account:        account,
				Password:       "",
				Description:    "",
				CubeSignerUser: "",
				OrganizationId: 1,
				Organization: &dto.RespOrg{
					ID:          1,
					Name:        organization,
					Description: "",
				},
			}, nil,
		)

	mockCSRepo.EXPECT().ToSign(ctx, userID, fmt.Sprintf("Key#Stellar_%s", from)).
		Times(1).
		Return(nil, assert.AnError)

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
	}
	_, err := svc.ToSign(ctx, &dto.ReqToSign{Account: account, Organization: organization, From: from})
	if assert.Error(t, err) {
		assert.Equal(t, assert.AnError, err)
	}
}

func TestToSignWithAccountNotFoundError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := repo.NewMockOAuth(ctrl)
	mockCSRepo := repo.NewMockCubeSigner(ctrl)

	ctx := context.TODO()

	from := "from"
	account := "v3n"
	organization := "epic"

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, &dto.ReqOrgAcn{
		OrgName: "epic",
		AcnName: "v3n",
	}).
		Times(1).
		Return(nil, assert.AnError)

	svc := &service{
		oauthRepo: mockOAuthRepo,
		csRepo:    mockCSRepo,
	}
	_, err := svc.ToSign(ctx, &dto.ReqToSign{Account: account, Organization: organization, From: from})
	if assert.Error(t, err) {
		assert.Equal(t, assert.AnError, err)
	}
}
