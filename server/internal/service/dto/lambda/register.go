package lambda

import (
	"io"

	dtoOrg "github.com/57blocks/auto-action/server/internal/service/dto/organization"
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
		_          struct{}
		Lambdas    []RespLamBrief `json:"lambdas"`
		Schedulers []RespSchBrief `json:"schedulers"`
	}

	RespLamBrief struct {
		_       struct{}
		Name    string `json:"function_name"`
		Arn     string `json:"function_arn"`
		Runtime string `json:"runtime"`
		Handler string `json:"handler"`
		Version string `json:"version"`
	}
	RespSchBrief struct {
		_              struct{}
		Arn            string `json:"schedule_arn"`
		BoundLambdaArn string `json:"bound_lambda_arn"`
	}
)
