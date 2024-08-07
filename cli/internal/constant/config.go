package constant

import (
	"fmt"
	"strconv"
)

type Config string

const (
	ConfigurationType Config = "toml"
	ConfigName        Config = ".st3llar"
	CredentialName    Config = ".st3llar-cred"

	EnvPrefix Config = "ST3LLAR"
)

func (cc Config) ValStr() string {
	return string(cc)
}

func (cc Config) ValInt() int {
	intVal, err := strconv.Atoi(cc.ValStr())
	if err != nil {
		fmt.Printf("converting <%s> to int error: %s\n", cc.ValStr(), err.Error())
		return 0
	}

	return intVal
}

func (cc Config) ValFloat32() float32 {
	floatValue, err := strconv.ParseFloat(cc.ValStr(), 32)
	if err != nil {
		fmt.Printf("Error converting <%s> to float32 error: %s\n", cc.ValStr(), err.Error())
		return float32(0.0)
	}

	return float32(floatValue)
}

func (cc Config) ValFloat64() float64 {
	floatValue, err := strconv.ParseFloat(cc.ValStr(), 64)
	if err != nil {
		fmt.Printf("Error converting <%s> to float64 error: %s\n", cc.ValStr(), err.Error())
		return 0.0
	}

	return floatValue
}
