# AutoAction


## AWS preparation
Terraform
1. VPC init with pub/pir subnets, and get the subnet ids.
2. Security group init for ALB, Application, RDS and public access.
3. Execution roles:
   a. Execution role for Lambda: CloudWatch logs, log groups, log streams and put events.
   b. Execution role for Scheduler: The `Resource` should involve **all** the Lambdas in the account.
   c. Execution role for ECS task: ecr and log related.
4. Secret Manager:
   1. Server secret key.


## DB migration
There is an initial version: `000000_init`, which aims at:
1. Solve the problem of the dirty version at the beginning. [Issue Ref](https://github.com/golang-migrate/migrate/issues/282#issuecomment-660760237)
2. The init version does nothing except: establish the changelog table: `schema_migrations`
3. If any error in migration which leads to a dirty version, fix migrations, then it will be re-executed when the 
   server starts.
4. If the fixed version is dirty still, go back to step `3`.
5. There exists some data migrations required:
   - Insert the VPC configuration when the Amazon infrastructure is ready.
     - Which subnets are going to use to host the BE endpoint.
     - Security groups.
   - Insert the organization in use.
   - Insert the initial user account.
   - CubeSigner related data.


## ESLint

Using `npm install eslint globals` to install ESLint for local environment.


## CubeSigner

### CubeSigner Session
1. Login with MFA:
   - `cs login -s google --session-lifetime 31536000 --auth-lifetime 600 --refresh-lifetime 31536000`
     - Command above will create a root/admin session with a quick-expired auth token, and a long-lived session and 
       refresh token.
     - The refresh frequency is based on the `--auth-lifetime 600`, which will keep the auth token alive and 
       refreshed short enough at the same time.
   - Response sample: 
        ```
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
   - Explanation:
     - `expiration`: `1756972127`, from the `--session-lifetime 31536000`, which is 12 months later.
     - `token`: if you use this token to sign, you will get an `ImproperSessionScope` error which indicates that 
       management token should not be used to sign.
     - `session_info`: 
       - This is like a change-log/change-version of the session refreshment, the `epoch` will be added 1 
         everytime you refreshed the session.
       - The session expiration will not be longer when refreshed.
   - Practice:
     1. Init a root/admin/management session with:
        1. long **session** lifetime, currently for 1 year.
        2. long **refresh** lifetime, currently, the same as the session lifetime.
        3. short **auth** lifetime, currently for 10 minutes.
     2. Hanging a recurring job to refresh the root session with a fixed interval to keep it alive, which is based 
        on the `--auth-lifetime`.
        1. We could make the fixed interval longer or shorter which is depends on your requirement.
        2. Using the token and required info to refresh, even the token is expired already.
     3. To log in again manually with MFA annually before the session and refresh expired.
     4. When a role/signer session is needed, we will use the root session to generate a **role**/**signer** session 
        with short lifetime for all scopes, which will be used to sign transactions. Generally, in some degree, we 
        could treat it as an one-time session.