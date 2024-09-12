package dto

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/aws/smithy-go/middleware"
)

type ReqURILambda struct {
	Lambda string `uri:"lambda" json:"lambda"`
}

// Info
type (
	RespInfo struct {
		_            struct{}
		ID           uint64     `json:"-"`
		AccountId    uint64     `json:"account_id"`
		FunctionName string     `json:"function_name"`
		FunctionArn  string     `json:"function_arn"`
		Runtime      string     `json:"runtime"`
		Role         string     `json:"role"`
		Handler      string     `json:"handler"`
		Description  string     `json:"description"`
		CodeSHA256   string     `json:"code_sha256"`
		Version      string     `json:"version"`
		RevisionID   string     `json:"revision_id"`
		Scheduler    Scheduler  `json:"scheduler" gorm:"foreignKey:lambda_id"`
		CreatedAt    *time.Time `json:"created_at"`
		UpdatedAt    *time.Time `json:"updated_at"`
	}

	Scheduler struct {
		LambdaID     uint64 `json:"-"`
		ScheduleName string `json:"schedule_name,omitempty"`
		ScheduleArn  string `json:"schedule_arn,omitempty"`
		Expression   string `json:"expression,omitempty"`
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
		_         struct{}
		Lambda    *RespLamBrief `json:"lambda"`
		Scheduler *RespSchBrief `json:"scheduler"`
	}

	RespLamBrief struct {
		_         struct{}
		AccountId uint64 `json:"account_id,omitempty"`
		Name      string `json:"function_name,omitempty"`
		Arn       string `json:"function_arn,omitempty"`
		Runtime   string `json:"runtime,omitempty"`
		Handler   string `json:"handler,omitempty"`
		Version   string `json:"version,omitempty"`
	}
	RespSchBrief struct {
		_              struct{}
		Arn            string `json:"schedule_arn,omitempty"`
		Name           string `json:"schedule_name,omitempty"`
		BoundLambdaArn string `json:"bound_lambda_arn,omitempty"`
	}
)

type (
	StdEventPayload struct {
		_            struct{}
		Organization string `json:"organization"`
		Account      string `json:"account"`
	}
)

// RespRemove related
type RespRemove struct {
	_         struct{}
	Lambdas   RespLamBrief `json:"lambda"`
	Scheduler RespSchBrief `json:"scheduler"`
}
