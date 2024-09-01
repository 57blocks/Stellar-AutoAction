package lambda

type (
	ReqLogs struct {
		LambdaName string `uri:"lambda_name" json:"lambda_name"`
	}

	RespLogs struct {
		Logs []string `json:"logs"`
	}
)
