package lambda

import (
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/aws/smithy-go/middleware"
)

type (
	ReqTrigger struct {
		Lambda  string `uri:"lambda"`
		Payload string `json:"payload"`
	}

	RespTrigger struct {
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
	RespTrgOpt func(respTrigger *RespTrigger)
)

func BuildRespTrigger(opts ...RespTrgOpt) *RespTrigger {
	respTrigger := &RespTrigger{}

	for _, opt := range opts {
		opt(respTrigger)
	}

	return respTrigger
}

func WithTriggerResp(resp *lambda.InvokeOutput) RespTrgOpt {
	return func(respTrigger *RespTrigger) {
		respTrigger.ExecutedVersion = resp.ExecutedVersion
		respTrigger.FunctionError = resp.FunctionError
		respTrigger.LogResult = resp.LogResult
		respTrigger.Payload = resp.Payload
		respTrigger.StatusCode = resp.StatusCode
		respTrigger.ResultMetadata = resp.ResultMetadata
	}
}
