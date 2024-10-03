# AutoAction Server

## AWS Infrastructure Preparation

Use Terraform to set up the following AWS resources:

1. VPC with public and private subnets
2. Security groups for ALB, Application, RDS, and public access
3. Execution roles:
   a. Lambda execution role: Permissions for CloudWatch logs, log groups, log streams, and put events
   b. Scheduler execution role: Include all Lambda functions in the account as resources
   c. ECS task execution role: Permissions for ECR and logging
4. Secret Manager: Store server secret key

## Database Migration

The initial migration version `000000_init` serves the following purposes:

1. Resolves the "dirty version" issue at the beginning ([Reference](https://github.com/golang-migrate/migrate/issues/282#issuecomment-660760237))
2. Creates the `schema_migrations` changelog table
3. Handles migration errors:
   - If a migration fails, fix the issue and re-run the migration on server start
   - If the fixed version is still dirty, repeat the process
4. Required data migrations:
   - VPC configuration (subnets for BE endpoint, security groups)
   - Organization data
   - Initial user account
   - CubeSigner-related data

## ESLint Setup

Install ESLint for the local environment:

<<<<<<< HEAD
```
=======
```bash
>>>>>>> staging
npm install eslint globals
```

## Local Development Setup

To run the AutoAction server locally, follow these steps:

1. Clone the repository:

<<<<<<< HEAD
   ```
   git clone https://github.com/57blocks/AutoAction.git
   cd server
   ```

2. Install dependencies:

   ```
   go mod download
   ```

3. Set up environment variables:

   - TODO: Add environment variable setup instructions

4. Start the local PostgreSQL database (if not using a remote database):

   Start the local PostgreSQL database:

   ```
   docker run -d -p 5432:5432 postgres
   ```

   Create a database with the name `autoaction`(as same as the `RDS_DATABASE` in the env variables).

5. Start the server:
   ```
   go run main.go
   ```
=======
```bash
git clone https://github.com/57blocks/AutoAction.git
cd server
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

The server code requires some environment variables to be set. Create a `.env` file, following the example in `.env.example`:

```bash
cp .env.example .env
```

Here’s an explanation of the variables in the `.env` file.

**Key-Related Environment Variables**

- RSA_PRIVATE_KEY: The RSA private key, generated using the RSA asymmetric encryption algorithm, and then base64 encoded. You can follow [the instructions in the Infrastructure documentation](../infrastructure/README.md) to generate it. This private key corresponds to the `public_key` in the CLI configuration file.
- JWT_PUBLIC_KEY: The JWT public key, generated using the RSA asymmetric encryption algorithm, and then base64 encoded. Follow [the instructions in the Infrastructure documentation](../infrastructure/README.md) to generate it.
- JWT_PRIVATE_KEY: The JWT private key, also generated via the RSA asymmetric encryption algorithm and base64 encoded. It corresponds to the `JWT_PUBLIC_KEY`.

**AWS-Related Environment Variables**

You first need to add trust relationships for your local AWS user in the ECS Task role in AWS.
Go to ECS Task role in AWS: `Amazon Elastic Container Service -> Task definitions -> Select Task definition -> Select the latest revision -> Click on Task role -> Click on Trust relationships`. Add the following content:

```json
"Statement": [
    // .......
    {
        "Effect": "Allow",
        "Principal": {
            "AWS": "arn:aws:iam::12345678790:user/user@org.com"
        },
        "Action": "sts:AssumeRole"
    }
]
```

Next, use the AWS CLI to obtain a session token for the ECS Task role. The token is valid for 1 hour. Run the following command, and input the ARN of your ECS Task role (for example, `arn:aws:iam::12345678790:role/ecsTaskExecutionRole`):

```bash
aws sts assume-role --role-arn arn:aws:iam::12345678790:role/ecsTaskExecutionRole --role-session-name ecs-role-session
{
    "Credentials": {
        "AccessKeyId": "ASIARZN5DFBXG4GIRU2S",
        "SecretAccessKey": "wD6t/AkEXJHoKaFxbCC4dS6tzp0QqN/XfA4WtyQK",
        "SessionToken": "IQoJb3JpZ2luX2VjELH//////////wEaCXVzLWVhc3QtMiJGMEQCIAwQeG0yce3TK4g5hSgz5JB6pJTKta+TaUkc5azmAzNKAiBud/sQc3zKCb6BjLUFc7c4uXZykRvRlA59DKasa3Na8SqmAgjq//////////8BEAQaDDEyMzM0MDAwNzUzNCIMg2rgrDnVE55PiGvhKvoB7Z4fuQyYiGVmzil0ihKtGav7nfTwWZCWkFSd3dmSzqbH3QbKyXcVm780NuzQmZrTVt2okEofJuo3hJT34BedFW6J81qVYDfkch0rXjVql2cB7ReYkQtGHAj4HmWQkTNAYOjsBxscewb4nq9JbgzWwWcZIYEFBh4wvyL8Q7zNLkR9yJkjYJUCXCx5X0yIrLUgcxlEU/hM9Lj0T/L2rdv7PXaciApUhQktJegsAvbh2Am/idWJ79dcNS5b4pRa/uZ6xYK2tnAGYstYs2b25EJK3n9C7A0zqRm7saLbdrP0J/4sF8b7wGM3bXrNK8ccvsYwYpJRMZCMwgJFozCFm8+3BjqeAZdwTAVRgs5yLQClrBCsSOuJqeM4kQjx+9WR3sWc7YsDTHMcAb22pI3BylwAVHcBP0trSKhgvcZuIB0csOE6p6bO54EmzPXxSkmaR0PpWqig63Zb1DfRhdGF8zt3v/ub0ddfOfk3xxgbr6ZwzbJKjG6D23CLj4OyX3y1O2pdVj7br6uJ0HwwgGTv//sb8qDDU3fg27i7sI0T0Md1eSaf",
        "Expiration": "2024-09-25T09:44:53+00:00"
    },
    // ......
}
```

Enter the `AccessKeyId`, `SecretAccessKey`, and `SessionToken` into the `.env` file as `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_SESSION_TOKEN`, respectively.

Here’s an example `.env` file with AWS variables:

```bash
AWS_REGION=us-east-2
AWS_ACCESS_KEY_ID=ASIARZN5DFBXG4GIRU2S
AWS_SECRET_ACCESS_KEY="wD6t/AkEXJHoKaFxbCC4dS6tzp0QqN/XfA4WtyQK"
AWS_SESSION_TOKEN="IQoJb3JpZ2luX2VjELH//////////wEaCXVzLWVhc3QtMiJGMEQCIAwQeG0yce3TK4g5hSgz5JB6pJTKta+TaUkc5azmAzNKAiBud/sQc3zKCb6BjLUFc7c4uXZykRvRlA59DKasa3Na8SqmAgjq//////////8BEAQaDDEyMzM0MDAwNzUzNCIMg2rgrDnVE55PiGvhKvoB7Z4fuQyYiGVmzil0ihKtGav7nfTwWZCWkFSd3dmSzqbH3QbKyXcVm780NuzQmZrTVt2okEofJuo3hJT34BedFW6J81qVYDfkch0rXjVql2cB7ReYkQtGHAj4HmWQkTNAYOjsBxscewb4nq9JbgzWwWcZIYEFBh4wvyL8Q7zNLkR9yJkjYJUCXCx5X0yIrLUgcxlEU/hM9Lj0T/L2rdv7PXaciApUhQktJegsAvbh2Am/idWJ79dcNS5b4pRa/uZ6xYK2tnAGYstYs2b25EJK3n9C7A0zqRm7saLbdrP0J/4sF8b7wGM3bXrNK8ccvsYwYpJRMZCMwgJFozCFm8+3BjqeAZdwTAVRgs5yLQClrBCsSOuJqeM4kQjx+9WR3sWc7YsDTHMcAb22pI3BylwAVHcBP0trSKhgvcZuIB0csOE6p6bO54EmzPXxSkmaR0PpWqig63Zb1DfRhdGF8zt3v/ub0ddfOfk3xxgbr6ZwzbJKjG6D23CLj4OyX3y1O2pdVj7br6uJ0HwwgGTv//sb8qDDU3fg27i7sI0T0Md1eSaf"
AWS_ECS_TASK_ROLE=arn:aws:iam::12345678790:role/ecsTaskExecutionRole
```

**Database-Related Environment Variables**

- DB_HOST: The database host, e.g., localhost.
- DB_PORT: The database port, e.g., 5432.
- DB_USER: The database username, e.g., postgres.
- DB_PASSWORD: The database password, e.g., password.
- DB_NAME: The database name, e.g., autoaction.

**Cube Signer-Related Environment Variables**

- CS_ORGANIZATION: The organization name in Cube Signer, typically prefixed with `Org#`, such as `Org#7bfdd921-bba7-505d-804d-36e2f2bf9357`.

4. Start the local PostgreSQL database (if not using a remote database):

```bash
docker run -d \
   --name postgres -p 5432:5432 \
   -e POSTGRES_PASSWORD=password \
   postgres
```

Create a database with the name `autoaction`(as same as the `RDS_DATABASE` in the env variables).

5. Start the server:

```bash
go run main.go
```
>>>>>>> staging

The server should now be running on `http://localhost:8080` (or the port specified in your configuration).

## CubeSigner Integration

### CubeSigner Session Management

1. Login with MFA:

<<<<<<< HEAD
   ```
=======
   ```bash
>>>>>>> staging
   cs login -s google --session-lifetime 31536000 --auth-lifetime 600 --refresh-lifetime 31536000
   ```

   - Creates a root/admin session with:
     - Quick-expiring auth token
     - Long-lived session and refresh token
   - Refresh frequency based on `--auth-lifetime 600`

2. Session management practice:
   a. Initialize root/admin session:
   - Long session lifetime (1 year)
   - Long refresh lifetime (same as session lifetime)
   - Short auth lifetime (10 minutes)
     b. Implement a recurring job to refresh the root session based on `--auth-lifetime`
     c. Manually re-authenticate with MFA annually before session expiration
     d. Generate short-lived role/signer sessions from the root session for transaction signing

### CubeSigner Session Details

Login response sample:

```json
{
  "org_id": "Org#...",
  "role_id": null,
  "expiration": 1756972127,
  "purpose": "OIDC-auth session with scopes [ManageAll]",
  "token": "3d6fd7397:...",
  "refresh_token": "3d6fd7397:...",
  "env": {
    "Dev-CubeSignerStack": {
      "ClientId": "1tiou9ecj058khiidmhj4ds4rj",
      "GoogleDeviceClientId": "59575607964-nc9hjnjka7jlb838jmg40qes4dtpsm6e.apps.googleusercontent.com",
      "GoogleDeviceClientSecret": "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",
      "Region": "us-east-1",
      "UserPoolId": "us-east-1_RU7HEslOW",
      "SignerApiRoot": "https://gamma.signer.cubist.dev",
      "DefaultCredentialRpId": "cubist.dev",
      "EncExportS3BucketName": null,
      "DeletedKeysS3BucketName": null
    }
  },
  "session_info": {
    "auth_token": "keCmhik9...",
    "auth_token_exp": 1725522527,
    "epoch": 1,
    "epoch_token": "LChINEm9si...",
    "refresh_token": "Z9OmjYwih...",
    "refresh_token_exp": 1728028127,
    "session_id": "7dfa49f5-xx-xx-xx-xx"
  }
}
```

Explanation:

- `expiration`: 1756972127 (12 months from session creation)
- `token`: Management token, not for signing (will result in `ImproperSessionScope` error if used for signing)
- `session_info`:
  - Tracks session refreshes (epoch increments with each refresh)
  - Session expiration remains unchanged after refreshes

Note: When signing transactions, use a short-lived role/signer session generated from the root session.
