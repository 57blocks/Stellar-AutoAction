package lambda

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type (
	ReqRegister struct {
		_     struct{}
		Files []io.Reader `json:"files"`
	}

	RespRegister struct {
		CFOs []*lambda.CreateFunctionOutput `json:"cfos"`
	}
)
