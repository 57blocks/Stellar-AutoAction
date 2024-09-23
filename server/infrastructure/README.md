# IaC

### Deploy environment

Currently, there is a sample environment: `staging` for development purposes.  
If any more environments are needed, create another directory besides the `staging`.  



### Variables injection

The variable file `variables.auto.tfvars` is ignored by git.  

1. Please add your own variables in both `variables.tf`, `variables.auto.tfvars` when you needed.
2. Do make sure that the variables in env, here: `staging/main.tf`, refer to the variables in modules, here: `modules/vpc/main.tf`.
    ```hcl
   // staging/main.tf left side is the input representation in modules
   // right side is the input representation in `staging`
   vpc_name = var.vpc_name
   // modules/vpc/main.tf
   azs = var.vpc_azs
    ```
3. Add default values to keep hcl validate. And also, it's recommended that keep the type and default value as it in 
   modules.

Here are some samples/notes for `variables.auto.tfvars` usage:  
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
  3. The key in env should be matched with references in code. Like: `JWT_PRIVATE_KEY`


### Complete Sample
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


### Apply

#### Preparation

1. AWS credentials
   1. credentials
   2. region
2. Terraform establish: [how](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
3. Sensitive Data
   1. Generate RSA key pairs for JWT
   2. Generate RSA key pairs for sensitive data encryption
   3. Make them above base64-encoded
4. Docker Image
   1. Build the image with proper arch

#### Terraform Apply

1. Comment the ECS module to apply the data required in it.
2. `terraform apply`
3. Push the image to ECR: [how](https://docs.aws.amazon.com/AmazonECR/latest/userguide/docker-push-ecr-image.html)
4. Uncomment the ECS module and `terraform apply` again.
