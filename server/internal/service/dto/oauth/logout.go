package oauth

type (
	ReqLogout struct {
		_     struct{}
		Token string `json:"token"`
	}

	RespLogout struct{}
)
