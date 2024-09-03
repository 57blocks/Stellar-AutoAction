# AutoAction

## AWS preparation
Terraform
1. VPC init with pub/pir subnets, and get the subnet ids.
2. Security group init for ALB, Application, RDS and public access.
3. Execution roles:
   a. Execution role for Lambda: CloudWatch logs, log groups, log streams and put events.
   b. Execution role for Scheduler: The `Resource` should involve **all** the Lambdas in the account.
   c. Execution role for ECS task: ecr and log related.


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