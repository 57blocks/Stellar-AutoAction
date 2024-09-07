package constant

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type Config string

const (
	ConfigurationType Config = "toml"
	ConfigName        Config = ".st3llar"
	CredentialName    Config = ".st3llar-credential"
)

func (cc Config) ValStr() string {
	return string(cc)
}

func (cc Config) ValInt() int {
	intVal, err := strconv.Atoi(cc.ValStr())
	if err != nil {
		slog.Error(fmt.Sprintf("converting <%s> to int error: %s", cc.ValStr(), err.Error()))
		os.Exit(1)
	}

	return intVal
}

func (cc Config) ValFloat32() float32 {
	floatValue, err := strconv.ParseFloat(cc.ValStr(), 32)
	if err != nil {
		fmt.Printf("Error converting <%s> to float32 error: %s", cc.ValStr(), err.Error())
		slog.Error(fmt.Sprintf("Error converting <%s> to float32 error: %s", cc.ValStr(), err.Error()))
		os.Exit(1)
	}

	return float32(floatValue)
}

func (cc Config) ValFloat64() float64 {
	floatValue, err := strconv.ParseFloat(cc.ValStr(), 64)
	if err != nil {
		slog.Error(fmt.Sprintf("Error converting <%s> to float64 error: %s", cc.ValStr(), err.Error()))
		os.Exit(1)
	}

	return floatValue
}
