package lambda

import (
	"context"
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

	mockRepo := repo.NewMockLambda(ctrl)

	ctx := context.TODO()
	request := &dto.ReqInfo{Lambda: "name/arn"}

	now := time.Now().UTC()
	sha256 := "E2X2ZXxocZcefGFb8lu2QnbYV8higgV2yYcJSwPLAA4="

	mockRepo.EXPECT().LambdaInfo(ctx, request).Times(1).
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
				Schedulers: []dto.Scheduler{
					{
						Expression:  "rate(1 minutes)",
						ScheduleArn: "arn:aws:scheduler:us-east-2:123340007534:schedule/default/epic-v3n-transfer",
					},
				},
				CreatedAt: &now,
				UpdatedAt: &now,
			}, nil,
		)

	cd := &service{
		lambdaRepo: mockRepo,
	}
	info, err := cd.Info(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, info.Schedulers, 1)
	assert.Equal(t, info.CodeSHA256, sha256)
}
