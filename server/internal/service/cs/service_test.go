package cs

import (
	"context"
	"testing"

	"github.com/57blocks/auto-action/server/internal/dto"
	model "github.com/57blocks/auto-action/server/internal/model/cs"
	"github.com/57blocks/auto-action/server/internal/model/oauth"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_ToSign_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := model.NewMockRepo(ctrl)

	ctx := context.Background()
	request := &dto.ReqToSign{Account: "v3n", Organization: "epic"}

	mockRepo.EXPECT().ToSign(ctx, request).Times(1).
		Return(
			[]*dto.RespToSign{{
				Organization: oauth.RespOrg{},
				Account:      oauth.RespUser{},
				Role:         "Role#_",
				Keys: []dto.RespCSKey{{
					Key:    "Key#_St3llar_",
					Scopes: []string{"scope#_1", "scope#_2"},
				}},
			}}, nil,
		)

	cd := &conductor{csRepo: mockRepo}
	roles, err := cd.ToSign(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, roles, 1)
	assert.Equal(t, "Role#_", roles[0].Role)
	assert.Equal(t, "Key#_St3llar_", roles[0].Keys[0].Key)
}

func Test_ToSign_With_Empty_Roles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := model.NewMockRepo(ctrl)

	ctx := context.Background()
	request := &dto.ReqToSign{Account: "v3n", Organization: "epic"}

	mockRepo.EXPECT().ToSign(ctx, request).Times(1).
		Return(
			[]*dto.RespToSign{}, nil,
		)

	cd := &conductor{csRepo: mockRepo}
	roles, err := cd.ToSign(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, roles, 0)
	assert.Empty(t, roles)
}

func Test_ToSign_With_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := model.NewMockRepo(ctrl)

	ctx := context.Background()
	request := &dto.ReqToSign{Account: "v3n", Organization: "epic"}

	mockRepo.EXPECT().ToSign(ctx, request).Times(1).
		Return(
			nil, assert.AnError,
		)

	cd := &conductor{csRepo: mockRepo}
	_, err := cd.ToSign(ctx, request)
	if assert.Error(t, err) {
		assert.Equal(t, assert.AnError, err)
	}
}
