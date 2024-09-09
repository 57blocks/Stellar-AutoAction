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

//func TestRegisterSuccess(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockRepo := repo.NewMockLambda(ctrl)
//	mockAmazon := amazonx.NewMockAmazon(ctrl)
//
//	ctx := context.TODO()
//
//	request := &http.Request{
//		Method: "POST",
//		URL:    nil,
//		Header: map[string][]string{
//			"Content-Type":  {"multipart/form-data"},
//			"Authorization": {"token"},
//		},
//		Body:             nil,
//		GetBody:          nil,
//		ContentLength:    0,
//		TransferEncoding: nil,
//		Close:            false,
//		Host:             "",
//		Form:             nil,
//		PostForm:         nil,
//		//MultipartForm: &multipart.Form{
//		//	Value: map[string][]string{"file"},
//		//	File:  map[string][]*multipart.FileHeader{{"file": []*multipart.FileHeader{fileHeader}}},
//		//},
//	}
//
//	mockAmazon.EXPECT().RegisterLambda(ctx, request).
//		Times(1).
//		Return(
//			&lambda.CreateFunctionOutput{
//				Architectures: []types.Architecture{
//					types.ArchitectureX8664,
//					types.ArchitectureArm64,
//				},
//				FunctionArn:  aws.String("FunctionArn"),
//				FunctionName: aws.String("FunctionName"),
//				CodeSha256:   aws.String("CodeSha256"),
//				CodeSize:     0,
//				Description:  aws.String("Description"),
//				Runtime:      types.RuntimeNodejs20x,
//				Version:      aws.String("$LATEST"),
//				Handler:      aws.String("handler.handler"),
//			}, nil,
//		)
//	//mockRepo.EXPECT().PersistRegResult(ctx, func(tx *gorm.DB) error { return nil }).
//	mockRepo.EXPECT().PersistRegResult(ctx, nil).
//		Times(1).
//		Return(nil)
//
//	cd := &conductor{
//		lambdaRepo: mockRepo,
//		amazon:     mockAmazon,
//	}
//	register, err := cd.Register(ctx, request)
//	fmt.Println(register)
//	assert.NoError(t, err)
//	//assert.Equal(t, register, "")
//}

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

	cd := &conductor{
		lambdaRepo: mockRepo,
	}
	info, err := cd.Info(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, info.Schedulers, 1)
	assert.Equal(t, info.CodeSHA256, sha256)
}
