package util

import (
	"fmt"
	"strings"

	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
)

// parse key_id(format: Key#Stellar_<address>) to get the address
func GetAddressFromCSKey(key string) (string, error) {
	if !strings.Contains(key, "_") {
		return "", errorx.Internal("invalid key")
	}
	return strings.Split(key, "_")[1], nil
}

// assemble the key_id(format: Key#Stellar_<address>) from address
func GetCSKeyFromAddress(address string) string {
	return fmt.Sprintf("Key#Stellar_%s", address)
}
