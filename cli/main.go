package main

import (
	"github.com/57blocks/auto-action/cli/internal/command"
	_ "github.com/57blocks/auto-action/cli/internal/command/action"
	_ "github.com/57blocks/auto-action/cli/internal/command/auth"
	_ "github.com/57blocks/auto-action/cli/internal/command/general"
	_ "github.com/57blocks/auto-action/cli/internal/command/wallet"
)

func main() {
	command.Execute()
}
