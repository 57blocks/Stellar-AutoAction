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

```
npm install eslint globals
```

## Local Development Setup

To run the AutoAction server locally, follow these steps:

1. Clone the repository:

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

The server should now be running on `http://localhost:8080` (or the port specified in your configuration).

## CubeSigner Integration

### CubeSigner Session Management

1. Login with MFA:

   ```
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
