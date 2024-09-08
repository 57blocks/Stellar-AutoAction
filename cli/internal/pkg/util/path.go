package util

import (
	"os"
	"strings"

	"github.com/57blocks/auto-action/cli/internal/constant"

	"github.com/spf13/cobra"
)

// Home $HOME path
func Home() string {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return home
}

// DefaultPath returns the default path of the config file
func DefaultPath() string {
	return Home() + "/" + constant.ConfigName.ValStr()
}

// DefaultCredPath returns the default path of the credential file
func DefaultCredPath() string {
	return Home() + "/" + constant.CredentialName.ValStr()
}

func ParseReqPath(input string) string {
	return strings.ReplaceAll(input, "#", "%23")
}
