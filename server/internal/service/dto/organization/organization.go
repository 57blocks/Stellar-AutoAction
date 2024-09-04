package organization

type (
	RespOrgRoleKey struct {
		_          struct{}
		CSRoleKeys []RespCSRoleKey `json:"cs_role_keys"`
	}

	RespCSRoleKey struct {
		_        struct{}
		CSRoleID string   `json:"role"`
		CSKeyID  string   `json:"key"`
		CSScopes []string `json:"scopes"`
	}
)

type (
	ReqSDKRequired struct {
		Organization string `json:"organization"`
		Account      string `json:"account"`
	}

	RespSDKRequired struct {
		Token        string          `json:"token"`
		Organization string          `json:"organization"`
		Keys         []RespCSRoleKey `json:"cs_role_keys"`
	}
)

type (
	RespRefreshRootSession struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		Expiration   int    `json:"expiration"`
		SessionInfo  struct {
			SessionID       string `json:"session_id"`
			AuthToken       string `json:"auth_token"`
			RefreshToken    string `json:"refresh_token"`
			Epoch           int    `json:"epoch"`
			EpochToken      string `json:"epoch_token"`
			AuthTokenExp    int    `json:"auth_token_exp"`
			RefreshTokenExp int    `json:"refresh_token_exp"`
		} `json:"session_info"`
	}
)
