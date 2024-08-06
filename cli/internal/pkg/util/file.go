package util

import "os"

func IsExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}
