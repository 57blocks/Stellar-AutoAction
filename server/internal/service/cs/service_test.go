package cs

import (
	"context"
	"testing"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/repo"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestToSignSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockCubeSigner(ctrl)

	ctx := context.TODO()
	request := &dto.ReqToSign{Account: "v3n", Organization: "epic"}

	mockRepo.EXPECT().ToSign(ctx, request).Times(1).
		Return(
			[]*dto.RespCSKey{{}}, nil,
		)

	cd := &service{csRepo: mockRepo}
	toSign, err := cd.ToSign(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, "Role#_", toSign)
}

func TestToSignWithEmptyRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockCubeSigner(ctrl)

	ctx := context.TODO()
	request := &dto.ReqToSign{Account: "v3n", Organization: "epic"}

	mockRepo.EXPECT().ToSign(ctx, request).Times(1).
		Return(
			[]*dto.RespCSKey{}, nil,
		)

	cd := &service{csRepo: mockRepo}
	roles, err := cd.ToSign(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, roles, 0)
	assert.Empty(t, roles)
}

func TestToSignWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockCubeSigner(ctrl)

	ctx := context.TODO()
	request := &dto.ReqToSign{Account: "v3n", Organization: "epic"}

	mockRepo.EXPECT().ToSign(ctx, request).Times(1).
		Return(
			nil, assert.AnError,
		)

	cd := &service{csRepo: mockRepo}
	_, err := cd.ToSign(ctx, request)
	if assert.Error(t, err) {
		assert.Equal(t, assert.AnError, err)
	}
}
