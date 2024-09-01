package lambda

type (
	ReqTrigger struct {
		Lambda string `uri:"lambda"`
	}

	RespTrigger struct {
		Status string `json:"status"`
	}
)
