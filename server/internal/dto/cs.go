package dto

// CubeSigner
type (
	RespAddCsKey struct {
		Keys []struct {
			KeyID string `json:"key_id"`
		} `json:"keys"`
	}
)
