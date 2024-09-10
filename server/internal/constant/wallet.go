package constant

type StellarNetworkTypeEnum struct {
	TestNet string
	MainNet string
}

var StellarNetworkType = StellarNetworkTypeEnum{
	TestNet: "testnet",
	MainNet: "mainnet",
}
