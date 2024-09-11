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

	ListWalletsRespInfo struct {
		Data []ListWalletRespInfo `json:"data"`
	}

	VerifyWalletReqInfo struct {
		Address string `uri:"address"`
	}

	VerifyWalletRespInfo struct {
		Address string `json:"address"`
		IsValid bool   `json:"is_valid"`
	}
)
