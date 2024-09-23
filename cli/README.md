# AutoAction CLI

AutoAction CLI Tool: `autoaction`

## Introduction

`autoaction` is the command-line interface (CLI) tool for AutoAction tasks, used to manage and execute various operations related to AutoAction.

## Installation

Please refer to the installation instructions in the project's root directory for installation.

## Main Commands

The `autoaction` CLI includes the following main command groups:

1. **oauth** - User authentication related operations
2. **wallet** - Wallet address management
3. **lambda** - Lambda function management
4. **general** - General CLI settings

Use `autoaction help` to view all available commands.

## Configuration

Use the `autoaction configure` command to configure the CLI. The main configuration items include:

- Credential file path
- API endpoint
- Log level
- Trace source settings

## Scheduler Expression Types

`autoaction` supports the following types of scheduler expressions:

1. **Rate-based**: [rate-based](https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#rate-based)
2. **Cron-based**: [cron-based](https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#cron-based)
3. **One-time**: [one-time](https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#one-time)

## Development

### Adding New Commands

1. **Group Commands by Responsibility**: All commands should be grouped according to their responsibilities.
2. **Create Subcommand Template**: Execute `cobra-cli add sub-command` in the `.workspace/cli/` directory to create a new subcommand template file.
3. **Initialize the Subcommand**: Remember to initialize the newly added subcommand package in the `main.go` file:

```go
// Example initialization in main.go
import (
    "path/to/new/subcommand"
)

func main() {
    // Existing initialization code
    rootCmd.AddCommand(subcommand.NewCmd())
    // More initialization code
}
```

Make sure to replace `"path/to/new/subcommand"` with the actual path to your new subcommand package and initialize it appropriately.