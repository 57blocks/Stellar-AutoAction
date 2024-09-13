package lambda

import (
	"github.com/57blocks/auto-action/server/internal/constant"
	"github.com/gin-gonic/gin"
	"testing"
	"time"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/repo"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInfoSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := repo.NewMockLambda(ctrl)
	mockOAuthRepo := repo.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	request := &dto.ReqURILambda{Lambda: "name/arn"}
	distinguisher := "name/arn"
	accountID := uint64(123)

	now := time.Now().UTC()
	sha256 := "E2X2ZXxocZcefGFb8lu2QnbYV8higgV2yYcJSwPLAA4="
	schARN := "arn:aws:scheduler:...:schedule/default/transfer"

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, distinguisher).Times(1).
		Return(
			&dto.RespInfo{

				CodeSHA256:   sha256,
				Description:  "",
				FunctionArn:  "arn:aws:lambda:us-east-2:123340007534:function:epic-v3n-transfer",
				FunctionName: "epic-v3n-transfer",
				Handler:      "transfer.handler",
				RevisionID:   "3c9a3513-5e43-419e-ae5b-aeeb459e44e3",
				Role:         "arn:aws:iam::123340007534:role/LambdaExecutionRole",
				Runtime:      "nodejs20.x",
				Version:      "$LATEST",
				Scheduler: dto.Scheduler{
					Expression:  "rate(1 minutes)",
					ScheduleArn: schARN,
				},
				CreatedAt: &now,
				UpdatedAt: &now,
			}, nil,
		)

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
	}
	info, err := cd.Info(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, info.CodeSHA256, sha256)
	assert.Equal(t, info.Scheduler.ScheduleArn, schARN)
}
