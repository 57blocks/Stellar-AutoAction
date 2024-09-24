package dto

// CubeSigner
type (
	RespAddCsKey struct {
		Keys []struct {
			KeyID string `json:"key_id"`
		} `json:"keys"`
	}

	RespAddCsRole struct {
		Name   string `json:"name"`
		RoleId string `json:"role_id"`
	}
)
