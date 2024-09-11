package dto

type (
	RespCreateWallet struct {
		Address string `json:"address"`
	}

	ReqRemoveWallet struct {
		Address string `uri:"address"`
	}

	RespListWallet struct {
		Address string `json:"address"`
	}

	RespListWallets struct {
		Data []RespListWallet `json:"data"`
	}

	ReqVerifyWallet struct {
		Address string `uri:"address"`
	}

	RespVerifyWallet struct {
		Address string `json:"address"`
		IsValid bool   `json:"is_valid"`
	}
)
