# Stellar AutoAction CLI

Stellar AutoAction CLI Tool: `autoaction`

## Introduction

`autoaction` is the command-line interface (CLI) tool for Stellar AutoAction tasks, used to manage and execute various operations related to Stellar AutoAction.

## Installation

You can install the CLI tool in two ways as below.

### Installation from Source

Build the client:

```bash
# Clone the repository
git clone https://github.com/57blocks/Stellar-AutoAction.git
cd Stellar-AutoAction

# Install client dependencies
cd Stellar-AutoAction/cli
go mod download

# Build the client
go build -o autoaction .
```

After the command executes, a `autoaction` executable file will be generated in the current directory. Running this file will launch the CLI tool.

PS: It is recommended to move the file to a directory in your system’s PATH, so it can be used from any location.

### Running from Source

Run the client:

```bash
cd Stellar-AutoAction/cli
go run main.go
```

## Configuring the CLI Tool

The first time you run the CLI tool, a configuration file will be automatically generated. This file is usually saved in `~/.autoaction.toml`. The default configuration is as follows:

```toml
[general]
  logx = "Info"
  source = "OFF"
  public_key = ""

[bound_with]
  credential = "/Users/{username}/.autoaction-credential"
  endpoint = ""
```

There are two configuration items that need manual setup:

- public_key should be set to the value of the Terraform environment variable `rsa_public_key`.
- endpoint should be the address of the ECS service exposed on AWS. You can find it on AWS: `EC2 -> Load Balancers -> select the relevant Load Balancer -> DNS name`. For example, if the DNS name is `autoaction-alb-365278.us-east-2.elb.amazonaws.com`, the endpoint value will be `http://autoaction-alb-365278.us-east-2.elb.amazonaws.com`.

## Main Commands

The `autoaction` CLI includes the following main command groups:

1. **auth** - User authentication related operations
2. **wallet** - Wallet address management
3. **action** - Action management
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
