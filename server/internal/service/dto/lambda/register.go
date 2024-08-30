package lambda

import (
	"io"

	dtoOrg "github.com/57blocks/auto-action/server/internal/service/dto/organization"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type (
	ReqSchedulerEvent struct {
		_               struct{}
		Token           string `json:"token"`
		dtoOrg.JWTPairs `json:"jwt_pairs"`
		CubeSignerPairs []dtoOrg.CubeSignerPairs `json:"cubesigner_pairs"`
	}
)

type (
	ReqRegister struct {
		_     struct{}
		Files []io.Reader `json:"files"`
	}

	RespRegister struct {
		_    struct{}
		CFOs []*lambda.CreateFunctionOutput `json:"cfos"`
	}
)
