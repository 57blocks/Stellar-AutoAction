package main

import (
	"github.com/57blocks/auto-action/cli/internal/command"
	_ "github.com/57blocks/auto-action/cli/internal/command/general"
	_ "github.com/57blocks/auto-action/cli/internal/command/lambda"
	_ "github.com/57blocks/auto-action/cli/internal/command/oauth"
	_ "github.com/57blocks/auto-action/cli/internal/command/wallet"
)

func main() {
	command.Execute()
}
