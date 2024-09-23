provider "aws" {
  region = var.region
}

data "aws_caller_identity" "current" {}

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
  ingress_rules       = ["http-80-tcp", "https-443-tcp"]

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
      rule                     = "all-all"
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
  ingress_rules       = ["postgresql-tcp"]

  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules       = ["all-all"]
}

// Load Balancer module
module "alb" {
  source = "./../modules/loadbalancer"

  alb_name            = var.alb_name
  alb_vpc_id          = module.vpc.vpc_id
  alb_subnets         = module.vpc.vpc_public_subnets
  alb_security_groups = [module.sg_alb.sg_id]
  alb_listener = {
    ex_http = {
      port     = 80
      protocol = "HTTP"

      forward = {
        target_group_key = "ecs_app"
      }
    }
  }

  alb_target_groups = {
    ecs_app = {
      name                              = "ecs-app"
      protocol                          = "HTTP"
      port                              = 8080
      target_type                       = "ip"
      deregistration_delay              = 5
      load_balancing_cross_zone_enabled = true
      protocol_version                  = "HTTP1"

      health_check = {
        enabled             = true
        healthy_threshold   = 5
        interval            = 300
        matcher             = "200"
        path                = "/up"
        port                = "traffic-port"
        protocol            = "HTTP"
        timeout             = 5
        unhealthy_threshold = 2
      }
      create_attachment = false
    }
  }
}

// ECR
module "ecr" {
  source = "./../modules/ecr"

  ecr_repository_name               = var.ecr_repository_name
  repository_read_write_access_arns = [module.ecs_task_execution_role.role_arn]
}

// IAM roles
module "scheduler_invocation_role" {
  source = "../modules/iam"

  role_name        = "auac_scheduler_invocation_role"
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
          "arn:aws:lambda:${var.region}:${data.aws_caller_identity.current.account_id}:function:*"
        ]
      }
    ]
  })
}

module "ecs_task_execution_role" {
  source = "../modules/iam"

  role_name        = "auac_ecs_task_execution_role"
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
          "ecr:BatchGetImage"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "logs:PutLogEvents",
          "logs:CreateLogStream",
          "logs:DescribeLogStreams",
          "logs:CreateLogGroup",
          "logs:DescribeLogGroups"
        ]
        Resource = "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:*"
      },
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = "arn:aws:secretsmanager:${var.region}:${data.aws_caller_identity.current.account_id}:secret:*"
      },
    ]
  })
}

// RDS
resource "aws_db_subnet_group" "default" {
  name       = "${var.vpc_name}-db-subnet-group"
  subnet_ids = module.vpc.vpc_public_subnets

  tags = {
    Name = "${var.vpc_name}-db-subnet-group"
  }
}

module "rds_password" {
  source = "./../modules/secretmanager"

  secret_name = var.rds_key_pairs
  secret_value = jsonencode({
    rds_username = var.rds_username
    rds_password = var.rds_password
  })
}

module "rds" {
  source = "./../modules/rds"

  depends_on = [module.rds_password, aws_db_subnet_group.default]

  rds_identifier = var.rds_identifier

  rds_db_name  = var.rds_db_name
  rds_username = jsondecode(module.rds_password.secret_value).rds_username
  rds_password = jsondecode(module.rds_password.secret_value).rds_password

  rds_vpc_security_group_ids = [module.sg_rds.sg_id]
  rds_db_subnet_group_name   = aws_db_subnet_group.default.name
  rds_subnet_ids             = aws_db_subnet_group.default.subnet_ids
}

// ECS
module "jwt_key_pairs" {
  source = "./../modules/secretmanager"

  secret_name = var.jwt_key_pairs
  secret_value = jsonencode({
    public_key  = var.jwt_public_key
    private_key = var.jwt_private_key
  })
}

module "rsa_key_pairs" {
  source = "./../modules/secretmanager"

  secret_name = var.rsa_key_pairs
  secret_value = jsonencode({
    public_key  = var.rsa_public_key
    private_key = var.rsa_private_key
  })
}

module "cs_key_pairs" {
  source = "./../modules/secretmanager"

  secret_name = var.cs_key_pairs
  secret_value = jsonencode({
    endpoint     = var.cs_endpoint
    organization = var.cs_organization
  })
}

module "ecs" {
  source = "./../modules/ecs"

  depends_on = [
    module.vpc,
    module.ecs_task_execution_role,
    module.alb,
    module.sg_ecs,
    module.jwt_key_pairs,
    module.ecr,
  ]

  ecs_cluster_name = var.ecs_cluster_name

  ecs_fargate_capacity_providers = var.ecs_fargate_capacity_providers

  ecs_services = {
    auac-service = {
      cpu    = 1024
      memory = 4096

      container_definitions = {
        auac-container = {
          cpu                = 512
          memory             = 1024
          essential          = true
          image              = nonsensitive("${module.ecr.ecr_repository_url}:latest")
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
              awslogs-group         = "/aws/ecs/ecs-services/auac-service"
              awslogs-region        = var.region
              awslogs-stream-prefix = "auto-actions"
              awslogs-create-group  = "true"
            }
          }
          environment = [
            {
              name  = "AWS_REGION"
              value = "us-west-2"
            },
            {
              name  = "BOUND_ENDPOINT"
              value = module.alb.alb_dns
            },
            {
              name  = "BOUND_NAME"
              value = "Horizon"
            },
            {
              name  = "LOG_LEVEL"
              value = "Info"
            },
            {
              name  = "LOG_ENCODING"
              value = "json"
            },
            {
              name  = "RDS_HOST"
              value = module.rds.db_address
            },
            {
              name  = "RDS_PORT"
              value = 5432
            },
            {
              name  = "RDS_DATABASE"
              value = "auac"
            },
            {
              name  = "RDS_SSLMODE"
              value = "require"
            },
            {
              name  = "WALLET_MAX"
              value = 10
            },
            {
              name  = "LAMBDA_MAX"
              value = 10
            },
            {
              name  = "AWS_ECS_TASK_ROLE"
              value = "arn:aws:iam::123340007534:role/AA-ECS-Task-Role" // TODO: tobe confirmed with Jimmy
            }
          ]
          secrets = [
            {
              name      = "RSA_PRIVATE_KEY"
              valueFrom = "${module.rsa_key_pairs.secret_arn}:private_key::"
            },
            {
              name      = "JWT_PUBLIC_KEY"
              valueFrom = "${module.jwt_key_pairs.secret_arn}:public_key::"
            },
            {
              name      = "JWT_PRIVATE_KEY"
              valueFrom = "${module.jwt_key_pairs.secret_arn}:private_key::"
            },
            {
              name      = "RDS_USER"
              valueFrom = "${module.rds_password.secret_arn}:rds_username::"
            },
            {
              name      = "RDS_PASSWORD"
              valueFrom = "${module.rds_password.secret_arn}:rds_password::"
            },
            {
              name      = "CS_ENDPOINT"
              valueFrom = "${module.cs_key_pairs.secret_arn}:endpoint::"
            },
            {
              name      = "CS_ORGANIZATION"
              valueFrom = "${module.cs_key_pairs.secret_arn}:organization::"
            }
          ]
        }
      }

      create_tasks_iam_role     = false
      create_task_exec_iam_role = false
      task_exec_iam_role_arn    = module.ecs_task_execution_role.role_arn

      subnet_ids = module.vpc.vpc_private_subnets

      load_balancer = {
        service = {
          target_group_arn = module.alb.target_groups
          container_name   = "auac-container"
          container_port   = 8080
        }
      }

      create_security_group = false
      security_group_ids    = [module.sg_ecs.sg_id]
    }
  }
}
