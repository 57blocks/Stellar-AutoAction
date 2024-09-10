package dto

type (
	CreateWalletRespInfo struct {
		Address string `json:"address"`
	}

	RemoveWalletReqInfo struct {
		Address string `uri:"address"`
	}
)
