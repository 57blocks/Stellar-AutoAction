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
