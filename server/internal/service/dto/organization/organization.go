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
	ReqKeys struct {
		Organization string `json:"organization"`
		Account      string `json:"account"`
	}
)
