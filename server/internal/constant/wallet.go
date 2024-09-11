package constant

type StellarNetworkTypeEnum struct {
	TestNet string
	MainNet string
}

var StellarNetworkType = StellarNetworkTypeEnum{
	TestNet: "Horizon-Testnet",
	MainNet: "Horizon",
}
