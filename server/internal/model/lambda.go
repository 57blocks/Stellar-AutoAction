package model

import "github.com/aws/aws-sdk-go-v2/service/lambda"

// Lambda model
type Lambda struct {
	ICU
	FunctionName string `json:"function_name"`
	FunctionArn  string `json:"function_arn"`
	Runtime      string `json:"runtime"`
	Timeout      uint8  `json:"timeout"`
	Role         string `json:"role"`
	Handler      string `json:"handler"`
	Description  string `json:"description"`
	CodeSHA256   string `json:"code_sha256"`
	Version      string `json:"version"`
	RevisionID   string `json:"revision_id"`
}

func (l *Lambda) TableName() string {
	return "lambda"
}

func (l *Lambda) TableNameWithAbbr() string {
	return "lambda AS l"
}

func TabNameLambda() string {
	return (&Lambda{}).TableName()
}

func TabNameLambdaAbbr() string {
	return (&Lambda{}).TableNameWithAbbr()
}

// LambdaScheduler model
type LambdaScheduler struct {
	ICU
	LambdaID    uint64 `json:"lambda_id"`
	ScheduleArn string `json:"schedule_arn"`
	Expression  string `json:"expression"`
}

func (l *LambdaScheduler) TableName() string {
	return "lambda_scheduler"
}

func (l *LambdaScheduler) TableNameWithAbbr() string {
	return "lambda_scheduler AS ls"
}

func TabNameLambdaSch() string {
	return (&LambdaScheduler{}).TableName()
}

func TabNameLambdaSchAbbr() string {
	return (&LambdaScheduler{}).TableNameWithAbbr()
}

// model builders and builder options
type (
	LambdaOpt    func(l *Lambda)
	SchedulerOpt func(l *LambdaScheduler)
)

// BuildLambda
// build the Lambda in optional pattern
func BuildLambda(opts ...LambdaOpt) *Lambda {
	l := new(Lambda)

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func WithLambdaResp(resp *lambda.CreateFunctionOutput) LambdaOpt {
	return func(l *Lambda) {
		l.FunctionName = *resp.FunctionName
		l.FunctionArn = *resp.FunctionArn
		l.Runtime = string(resp.Runtime)
		l.Timeout = uint8(*resp.Timeout)
		l.Role = *resp.Role
		l.Handler = *resp.Handler
		l.Description = *resp.Description
		l.CodeSHA256 = *resp.CodeSha256
		l.Version = *resp.Version
		l.RevisionID = *resp.RevisionId
	}
}

// BuildScheduler
// build the LambdaScheduler bound with Lambda in optional pattern
func BuildScheduler(opts ...SchedulerOpt) *LambdaScheduler {
	ls := new(LambdaScheduler)

	for _, opt := range opts {
		opt(ls)
	}

	return ls
}

func WithSchArn(arn string) SchedulerOpt {
	return func(l *LambdaScheduler) {
		l.ScheduleArn = arn
	}
}

func WithExpression(expression string) SchedulerOpt {
	return func(l *LambdaScheduler) {
		l.Expression = expression
	}
}

func WithLambdaID(lamIDRef uint64) SchedulerOpt {
	return func(l *LambdaScheduler) {
		l.LambdaID = lamIDRef
	}
}
