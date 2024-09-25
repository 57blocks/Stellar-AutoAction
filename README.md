# AutoAction

AutoAction is a comprehensive solution for managing OAuth, wallets, and tasks through a command-line interface (CLI) and a server endpoint.

## Command Line Interface (CLI)

### Goals

1. Provide a user-friendly CLI for AutoAction's OAuth, wallet, and task management functionalities.
2. Enhance the ease of creating and executing handler functions.

### Key Technologies

The CLI is developed using the following Go packages:

- `golang.org/x/crypto`: For cryptographic operations
- `github.com/spf13/cobra`: For building powerful modern CLI applications
- `github.com/spf13/viper`: For configuration solution
- `github.com/spf13/pflag`: For flag parsing
- `github.com/go-resty/resty/v2`: For HTTP and REST client operations

## Server Endpoint

The server component of AutoAction is built with a robust tech stack to ensure scalability, security, and performance.

### Key Features and Technologies

1. **Database**: Utilizes [PostgreSQL](https://www.postgresql.org/) for reliable data storage and management.
2. **API Framework**: Implements RESTful APIs using the [Gin](https://github.com/gin-gonic/gin) web framework, known for its performance and productivity.
3. **Cloud Services**: Leverages AWS services, including:
   - Lambda for serverless compute
   - ECS (Elastic Container Service) for container orchestration
   - EventBridge Scheduler for event-driven architecture
4. **Authentication**: Implements OAuth-based authentication for secure access control.
5. **Cryptography**: Integrates CubeSigner for robust cryptographic operations, ensuring the highest level of security for sensitive data.

## Getting Started

For detailed instructions on setting up and using AutoAction, please refer to the following documentation:

- [Init AWS environment](infrastructure/README.md)
- [CLI Setup and Usage Guide](cli/README.md)
- [Server Setup and Configuration](server/README.md)

## Contributing

We welcome contributions to AutoAction! Please read our [Contributing Guidelines](CONTRIBUTING.md) for more information on how to get started.

## License

AutoAction is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for more details.
