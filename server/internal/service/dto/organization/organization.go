package organization

type (
	RespRelatedRoleKey struct {
		_          struct{}
		CSRoleKeys []RespCSRoleKey `json:"cs_role_keys"`
	}

	RespCSRoleKey struct {
		_        struct{}
		CSRoleID string   `json:"cs_role_id"`
		CSKeyID  string   `json:"cs_key_id"`
		CSScopes []string `json:"cs_scopes"`
	}
)

type (
	RespRootSessionRefresh struct {
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
