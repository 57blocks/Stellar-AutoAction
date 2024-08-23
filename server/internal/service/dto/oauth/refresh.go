package oauth

type (
	ReqRefresh struct {
		_       struct{}
		Refresh string `json:"refresh"`
	}
)
