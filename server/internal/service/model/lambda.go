package model

type Lambda struct {
	ICU
	FunctionName string `json:"function_name"`
	FunctionARN  string `json:"function_arn"`
	Runtime      string `json:"runtime"`
	Role         string `json:"role"`
	Handler      string `json:"handler"`
	Description  string `json:"description"`
	CodeSHA256   string `json:"code_sha256"`
	Version      string `json:"version"`
	RevisionID   string `json:"revision_id"`
	VpcBound     bool   `json:"vpc_bound"`
	VpcID        string `json:"vpc_id"`
}

func (l *Lambda) TableName() string {
	return "lambda"
}

func (l *Lambda) TableNameWithAbbr() string {
	return "lambda AS l"
}

type LambdaVPC struct {
	ICU
	FunctionName string `json:"function_name"`
	FunctionARN  string `json:"function_arn"`
	Runtime      string `json:"runtime"`
	Role         string `json:"role"`
	Handler      string `json:"handler"`
	Description  string `json:"description"`
	CodeSHA256   string `json:"code_sha256"`
	Version      string `json:"version"`
	RevisionID   string `json:"revision_id"`
	VpcBound     bool   `json:"vpc_bound"`
	VpcID        string `json:"vpc_id"`
}

func (lv *LambdaVPC) TableName() string {
	return "lambda_vpc"
}

func (lv *LambdaVPC) TableNameWithAbbr() string {
	return "lambda_vpc AS lv"
}
