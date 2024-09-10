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
)
