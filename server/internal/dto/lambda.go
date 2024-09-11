package dto

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/aws/smithy-go/middleware"
)

// Info
type (
	ReqInfo struct {
		Lambda string `uri:"lambda" json:"lambda"`
	}

	RespInfo struct {
		_            struct{}
		ID           uint64      `json:"-"`
		FunctionName string      `json:"function_name"`
		FunctionArn  string      `json:"function_arn"`
		Runtime      string      `json:"runtime"`
		Role         string      `json:"role"`
		Handler      string      `json:"handler"`
		Description  string      `json:"description"`
		CodeSHA256   string      `json:"code_sha256"`
		Version      string      `json:"version"`
		RevisionID   string      `json:"revision_id"`
		Schedulers   []Scheduler `json:"schedulers" gorm:"foreignKey:lambda_id"`
		CreatedAt    *time.Time  `json:"created_at"`
		UpdatedAt    *time.Time  `json:"updated_at"`
	}

	Scheduler struct {
		LambdaID    uint64 `json:"-"`
		ScheduleArn string `json:"schedule_arn"`
		Expression  string `json:"expression"`
	}
)

// Invoke
type (
	ReqInvoke struct {
		Lambda  string `uri:"lambda"`
		Payload string `json:"payload"`
	}

	RespInvoke struct {
		_               struct{}
		ExecutedVersion *string

		// If present, indicates that an error occurred during function execution. Details
		// about the error are included in the response payload.
		FunctionError *string

		// The last 4 KB of the execution log, which is base64-encoded.
		LogResult *string

		// The response from the function, or an error object.
		Payload []byte

		// The HTTP status code is in the 200 range for a successful request. For the
		// RequestResponse invocation type, this status code is 200. For the Event
		// invocation type, this status code is 202. For the DryRun invocation type, the
		// status code is 204.
		StatusCode int32

		// Metadata pertaining to the operation's result.
		ResultMetadata middleware.Metadata
	}
	RespInvokeOpt func(respTrigger *RespInvoke)
)

func BuildRespInvoke(opts ...RespInvokeOpt) *RespInvoke {
	respTrigger := &RespInvoke{}

	for _, opt := range opts {
		opt(respTrigger)
	}

	return respTrigger
}

func WithInvokeResp(resp *lambda.InvokeOutput) RespInvokeOpt {
	return func(respTrigger *RespInvoke) {
		respTrigger.ExecutedVersion = resp.ExecutedVersion
		respTrigger.FunctionError = resp.FunctionError
		respTrigger.LogResult = resp.LogResult
		respTrigger.Payload = resp.Payload
		respTrigger.StatusCode = resp.StatusCode
		respTrigger.ResultMetadata = resp.ResultMetadata
	}
}

// ReqLogs cloudwatch event logs
type (
	ReqLogs struct {
		Lambda string `uri:"lambda"`
	}
)

// Register related
type (
	ReqRegister struct {
		_          struct{}
		Expression string
		Files      []*ReqFile
	}
	ReqFile struct {
		_     struct{}
		Name  string
		Bytes []byte
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

type (
	StdEventPayload struct {
		_            struct{}
		Organization string `json:"organization"`
		Account      string `json:"account"`
	}
)
