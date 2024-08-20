# AutoAction

## DB migration
There is an initial version: `000000_init`, which aims at:
1. Solve the problem of the dirty version at the version at the beginning. [Issue Ref](https://github.com/golang-migrate/migrate/issues/282#issuecomment-660760237)
2. The init version does nothing except: setup the tracing table: `schema_migrations`
3. If any error in migration which leads to a dirty version, fix it, and then it will be re-executed when the 
   server started.
4. If the fixed version is dirty still, go back to step `3`.