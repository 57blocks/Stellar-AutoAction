package util

import (
	"strings"
)

// parse key_id(format: Key#Stellar_<address>) to get the address
func GetAddressFromCSKey(key string) string {
	return strings.Split(key, "_")[1]
}
