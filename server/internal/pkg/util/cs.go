package util

import (
	"fmt"
	"strings"
)

// parse key_id(format: Key#Stellar_<address>) to get the address
func GetAddressFromCSKey(key string) string {
	return strings.Split(key, "_")[1]
}

// GetCSKeyFromAddress get the key_id from address
func GetCSKeyFromAddress(address string) string {
	return fmt.Sprintf("Key#Stellar_%s", address)
}
