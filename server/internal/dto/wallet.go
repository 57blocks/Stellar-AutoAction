package dto

type (
	CreateWalletRespInfo struct {
		Address string `json:"address"`
	}

	RemoveWalletReqInfo struct {
		Address string `uri:"address"`
	}

	ListWalletRespInfo struct {
		Address string `json:"address"`
	}

	ListWalletsResponse struct {
		Data []ListWalletRespInfo `json:"data"`
	}

	VerifyWalletReqInfo struct {
		Address string `uri:"address"`
		Env     string `json:"env"`
	}

	VerifyWalletRespInfo struct {
		Address string `json:"address"`
		IsValid bool   `json:"is_valid"`
	}
)
