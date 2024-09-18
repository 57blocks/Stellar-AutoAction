provider "aws" {
  region = var.region
}

// VPC module
module "vpc" {
  source = "./../modules/vpc"

  vpc_name        = var.vpc_name
  vpc_cidr        = var.vpc_cidr
  vpc_azs         = var.vpc_azs
  vpc_pub_subnets = var.vpc_pub_subnets
  vpc_pri_subnets = var.vpc_pri_subnets
}

// SG modules
module "sg_alb" {
  source = "./../modules/securitygroup"

  sg_vpc_id      = module.vpc.vpc_id
  sg_name        = var.sg_alb_name
  sg_description = "Security group for ALB"

  ingress_cidr_blocks = ["0.0.0.0/0"]
  ingress_rules       = ["http-80-tcp", "https-443-tcp"] // why two 80s in the rule?

  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules       = ["all-all"]
}

module "sg_ecs" {
  source = "./../modules/securitygroup"

  sg_vpc_id      = module.vpc.vpc_id
  sg_name        = var.sg_ecs_name
  sg_description = "Security group for Applications in ECS"

  ingress_with_source_security_group_id = [
    {
      description              = "http from ALB"
      rule                     = "http-80-tcp"
      source_security_group_id = module.sg_alb.sg_id
    }
  ]

  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules       = ["all-all"]
}

module "sg_rds" {
  source = "./../modules/securitygroup"

  sg_vpc_id      = module.vpc.vpc_id
  sg_name        = var.sg_rds_name
  sg_description = "Security group for RDS"

  ingress_cidr_blocks = ["0.0.0.0/0"]
  ingress_rules       = ["postgresql-tcp"] // why two 80s in the rule?

  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules       = ["all-all"]
}

// Load Balancer module
module "alb" {
  source = "./../modules/loadbalancer"

  alb_name            = var.alb_name
  alb_subnets         = module.vpc.vpc_public_subnets
  alb_security_groups = [module.sg_alb.sg_id] // List?
  alb_listener        = var.alb_listener
  alb_target_groups   = var.alb_target_groups
}

// IAM roles
module "scheduler_execution_role" {
  source = "../modules/iam"

  role_name        = "scheduler_execution_role"
  role_description = "Execution role for EventBridge Scheduler to invoke Lambda function"

  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        Effect = "Allow"
        Principal = {
          Service = ["scheduler.amazonaws.com"]
        }
        Action = ["sts:AssumeRole"]
      }
    ]
  })

  role_policy_name = "lambda_execution_role_policy"
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Action" : [
          "lambda:InvokeFunction"
        ],
        "Resource" : [
          // TODO: dynamic subject?
          "arn:aws:lambda:us-east-2:123340007534:function:*"
        ]
      }
    ]
  })
}

module "ecs_execution_role" {
  source = "../modules/iam"

  role_name        = "ecs_execution_role"
  role_description = "Execution role for ECS tasks"

  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        Effect = "Allow"
        Principal = {
          Service = ["ecs-tasks.amazonaws.com"]
        }
        Action = ["sts:AssumeRole"]
      }
    ]
  })

  role_policy_name = "ecs_execution_role_policy"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogStreams",
          "logs:CreateLogGroup",
          "logs:StartQuery",
          "logs:StopQuery",
          "logs:GetQueryResults",
          "logs:DescribeLogGroups"
        ]
        Resource = "*"
      },
    ]
  })
}

// Secrets Manager
module "jwt_private_key" {
  source = "./../modules/secretmanager"

  secret_name  = var.jwt_private_key
  secret_value = jsonencode({
    jwt_private_key = var.jwt_private_key_value
  })
}

module "jwt_public_key" {
  source = "./../modules/secretmanager"

  secret_name  = var.jwt_public_key
  secret_value = jsonencode({
    jwt_public_key = var.jwt_public_key_value
  })
}

module "rds_password" {
  source = "./../modules/secretmanager"

  secret_name  = var.rds_password
  secret_value = jsonencode({
    rds_password = var.rds_password_value
  })
}

//
# 创建一个秘密
resource "aws_secretsmanager_secret" "rds_password_secret" {
  name = "my-rds-password-secret"
}

# 设置秘密的值
resource "aws_secretsmanager_secret_version" "rds_password_version" {
  secret_id = aws_secretsmanager_secret.rds_password_secret.id
  secret_string = jsonencode({
    password = var.rds_password
  })
}

// RDS
module "rds" {
  source = "./../modules/rds"

  rds_identifier = var.rds_identifier

  rds_db_name  = var.rds_db_name
  rds_username = var.rds_username
  rds_password = module.rds_password.secret_value

  rds_db_subnet_group_name   = module.vpc.vpc_database_subnet_group_name
  rds_vpc_security_group_ids = module.vpc.vpc_private_subnets

  #   vpc_id             = module.vpc.vpc_id
  #   subnet_ids         = module.vpc.vpc_private_subnets
  #   security_group_ids = [module.sg_rds.sg_id]
}

// ECR
module "ecr" {
  source = "./../modules/ecr"

  ecr_repository_name               = var.ecr_repository_name
  repository_read_write_access_arns = [module.ecs_execution_role.role_arn]
}

// ECS module
module "ecs" {
  source = "./../modules/ecs"

  ecs_cluster_name = var.ecs_cluster_name

  ecs_fargate_capacity_providers = var.ecs_fargate_capacity_providers

  ecs_services = {
    aa-service = {
      cpu    = 1024
      memory = 4096

      # Container definition(s)
      container_definitions = {
        aa-service = {
          cpu       = 512
          memory    = 1024
          essential = true
          image     = nonsensitive(module.ecr.ecr_repository_arn) // TODO: image tag?
          firelens_configuration = {
            type = "fluentbit"
          }
          memory_reservation = 50
          port_mappings = [{
            containerPort = 8080
            hostPort      = 8080
            protocol      = "tcp"
            appProtocol   = "http"
          }]
          readonly_root_filesystem = false
          log_configuration = {
            logDriver = "awslogs"
            options = {
              awslogs-group         = "/aws/ecs/ecs-services/aa-service"
              awslogs-region        = var.region
              awslogs-stream-prefix = "auto-actions"
              awslogs-create-group  = "true"
            }
          }
          environment = {
            JWT_PRIVATE_KEY = nonsensitive(module.jwt_private_key.secret_arn)
            JWT_PUBLIC_KEY  = nonsensitive(module.jwt_public_key.secret_arn)
          }
        }
      }

      create_tasks_iam_role     = false
      create_task_exec_iam_role = false
      task_exec_iam_role_arn    = module.ecs_execution_role.role_arn

      vpc_id      = module.vpc.vpc_id
      vpc_subnets = module.vpc.vpc_public_subnets
      subnet_ids  = module.vpc.vpc_public_subnets

      load_balancer = {
        service = {
          target_group_arn = module.alb.target_groups
          container_name   = "aa-service"
          container_port   = 8080
        }
      }

      create_security_group = false
      security_group_ids    = [module.sg_ecs.sg_id]

      // TODO: 网络配置-服务角色: AWSServiceRoleForECS ?
    }
  }
}
