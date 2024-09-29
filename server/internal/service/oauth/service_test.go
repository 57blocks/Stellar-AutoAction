package oauth

import (
	"errors"
	"testing"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/repo"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
