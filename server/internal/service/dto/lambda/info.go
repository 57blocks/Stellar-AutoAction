package lambda

type (
	ReqInfo struct {
		Lambda string `uri:"lambda"`
	}

	RespInfo struct {
		ID           uint64      `json:"id"`
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
	}

	Scheduler struct {
		LambdaID    uint64 `json:"lambda_id"`
		ScheduleArn string `json:"schedule_arn"`
		Expression  string `json:"expression"`
	}
)
