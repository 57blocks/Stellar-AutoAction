# IaC

### 1. Deploy environment

Currently, there is a sample environment: `staging` for development purposes.  
If any more environments are needed, create another directory besides the `staging`.  
And almost of the configurations are the same as the `staging`, except for the `variables.auto.tfvars`, which is suitable for the new environment.


### 2. Paths
`.`: the root path of the project.  
`./infrastructure`: the path of the IaC.  
`./server`: the path of the server.  
`./cli`: the path of the CLI.  


### 3. Variables injection

The variable file `variables.auto.tfvars` is local-build.  

1. Add your own variables in both `variables.tf`, `variables.auto.tfvars` when new variables are needed.
2. Do make sure that the variables in env, is referred in modules.
    ```hcl
   // ./infrastructure/staging/main.tf
   // left side of `=`, is the input representation in modules. Here is `./modules/vpc`
   // right side of `=`, is the input representation in `./staging/main.tf`
   vpc_name = var.vpc_name
   
   // ./infrastructure/modules/vpc/main.tf
   name = var.vpc_name
    ```
3. Add default values to keep hcl validate. And also, it's recommended that keep the type and default value as it in 
   modules.

Here are some samples for `variables.auto.tfvars` usage:  
- General
    ```hcl
    region = "us-west-2"
    env    = "staging"
    ```
- VPC
    ```hcl
    vpc_name        = "aa-staging"
    vpc_cidr        = "172.16.0.0/24"
    vpc_azs         = ["us-west-2a", "..."]
    vpc_pub_subnets = ["172.16.0.0/28", "..."]
    vpc_pri_subnets = ["172.16.0.128/28", "..."]
    ```
    [VPC CIDR blocks](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-cidr-blocks.html)  
    [VPC subnet CIDR blocks](https://docs.aws.amazon.com/vpc/latest/userguide/subnet-sizing.html)  

- SecretsManager Key Pairs

    As you could see in the `staging/main.tf`:
    ```hcl
    module "rds_password" {
      source = "./../modules/secretmanager"
    
      secret_name = var.rds_key_pairs
      secret_value = jsonencode({
        rds_username = var.rds_username
        rds_password = var.rds_password
      })
    }
    ```
    Gathering related secrets together under a single secret name is a good practice.  
    
    **_Note_**: The JWT key paris are base64-encoded string, which is aiming at keeping the RSA PEM format.

- ECS  

  1. The cluster name must consist of alphanumerics, hyphens, and underscores.
  2. The Farget capacity provider is fixed in `variables.auto.tfvars` now:
     ```hcl
     ecs_fargate_capacity_providers = {
       FARGATE = {
         default_capacity_provider_strategy = {
           weight = 50
           base   = 20
         }
       }
       FARGATE_SPOT = {
         default_capacity_provider_strategy = {
           weight = 50
         }
       }
     }
     ```
  3. The key in env/secrets should be matched with references in code. Like: `JWT_PRIVATE_KEY`
  4. As for the `module "rsa_key_pairs"`, Generates a pair of RSA asymmetric encryption keys. 
     1. The public_key will be base64-encoded and added to your local config for the CLI.
     2. The private_key will be base64-encoded and stored in the SecretsManager for ECS service to use.
     3. Generation:
        1. Command or Online Tool
           ```shell
              openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
              openssl rsa -pubout -in private_key.pem -out public_key.pem
           ```
        2. Base64 Encode or Online Tool
           ```shell
              cat private_key.pem | base64
              cat public_key.pem | base64
           ```


### 4. Complete Sample
```hcl
region = "us-west-2"
env    = "staging"

// VPC
vpc_name        = "auac-staging"
vpc_cidr        = "172.64.0.0/24"
vpc_azs         = ["us-west-2a", "us-west-2b", "us-west-2c"]
vpc_pub_subnets = ["172.64.0.0/28", "172.64.0.16/28", "172.64.0.32/28"]
vpc_pri_subnets = ["172.64.0.128/28", "172.64.0.144/28", "172.64.0.160/28"]

// Security Groups
sg_alb_name = "auac-staging-alb-sg"
sg_ecs_name = "auac-staging-ecs-sg"
sg_rds_name = "auac-staging-rds-sg"

// ALB
alb_name = "auac-staging-alb"

// ECR
ecr_repository_name = "auto-actions"

// RDS SecretsManager
rds_key_pairs = "auac_rds_key_pairs"
rds_username  = "RDS_USERNAME"
rds_password  = "RDS_PASSWORD"

// RDS
rds_identifier = "auac-staging-postgres"
rds_db_name    = "auac"

// ECS SecretsManager
jwt_key_pairs   = "auac_jwt_key_pairs"
jwt_private_key = "LS0tLS1CRUdJTiBQUklW..."
jwt_public_key  = "LS0tLS1CRUdJTiBQVUJM..."

rsa_key_pairs   = "auac_rsa_key_pairs"
rsa_private_key = "LS0tLS1CRUdJTiBSU0Eg..."
rsa_public_key  = "LS0tLS1CRUdJTiBQVUJM..."

// ECS
ecs_cluster_name = "auac-staging"
ecs_fargate_capacity_providers = {
  FARGATE = {
    default_capacity_provider_strategy = {
      weight = 50
      base   = 20
    }
  }
  FARGATE_SPOT = {
    default_capacity_provider_strategy = {
      weight = 50
    }
  }
}
```


### 5. Apply

#### Preparation

1. AWS credentials
   1. credentials
   2. region
2. Terraform establish: [How](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
3. Terraform initialization: 
    ```shell
    // ./infrastructure/staging
    terraform init
    ```
4. Sensitive Data
   1. Generate RSA key pairs for JWT
   2. Generate RSA key pairs for sensitive data encryption
   3. Make them above base64-encoded
5. Docker Image 
   Here is the arch of `linux_amd64`:
    ```shell
    // `./server` 
    docker build --platform linux_amd64 -f ./build/Dockerfile -t auac:latest .
    ```

#### Terraform Apply

1. Terraform apply except for the ECS module.
    ```shell
    // ./infrastructure/staging
    terraform apply \
        -target=module.vpc \
        -target=module.sg_alb \
        -target=module.sg_ecs \
        -target=module.sg_rds \
        -target=module.alb \
        -target=module.ecr \
        -target=module.scheduler_invocation_role \
        -target=module.ecs_task_role \
        -target=module.ecs_task_execution_role \
        -target=aws_db_subnet_group.default \
        -target=module.rds_password \
        -target=module.rds \
        -target=module.jwt_key_pairs \
        -target=module.rsa_key_pairs \
        -target=module.cs_key_pairs
    ```
2. Push the image to ECR: [How](https://docs.aws.amazon.com/AmazonECR/latest/userguide/docker-push-ecr-image.html)
3. Apply for the ECS module.
    ```shell
    // ./infrastructure/staging
    terraform apply -target=module.ecs
    ```

#### Apply New Updates with ECS update

For example, I need to add another env/secrets for the ECS service:
1. Add the env into config in server side, then build into image and push to ECR.
2. Update the `./infrastructure/staging/main.tf`, to add the env vars.
3. Check the plan:
    ```shell
    // ./infrastructure/staging
    terraform plan
    ```
4. Apply the changes


#### Update the code only
                       
If the updates are not involved in the infrastructure.

1. Using the `update.sh`
   ```shell
   chmod +x ./update.sh
   ./update.sh
   ```
2. Alternatively, through the AWS CLI step by step according to the shell

3. Finally, you could do it in the Amazon console as well.