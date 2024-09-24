package util

import (
	"fmt"
	"strings"
)

// parse key_id(format: Key#Stellar_<address>) to get the address
func GetAddressFromCSKey(key string) string {
	return strings.Split(key, "_")[1]
}

// assemble the key_id(format: Key#Stellar_<address>) from address
func GetCSKeyFromAddress(address string) string {
	return fmt.Sprintf("Key#Stellar_%s", address)
}
