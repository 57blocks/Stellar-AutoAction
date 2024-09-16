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
          "arn:aws:lambda:us-east-2:123340007534:function:*" // TODO: dynamic subject
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

// ECS module
module "ecs" {
  source = "./../modules/ecs"

  cluster_name = var.ecs_cluster_name
  services = {
    aa-service = {
      cpu    = 512
      memory = 1024

      container_definitions = {
        aa-cli-server = {
          name      = "cli-server"
          cpu       = 512
          memory    = 1024
          essential = true
          image     = var.user_service_image
          port_mappings = [
            {
              name          = "aa-cli-server"
              containerPort = 8080
              hostPort      = 8080
              protocol      = "tcp"
              appProtocol   = "http"
            }
          ]
          readonly_root_filesystem = false
          log_configuration = {
            logDriver = "awslogs"
            options = {
              awslogs-group         = "/aws/ecs/aa-service/aa-cli-server"
              awslogs-region        = var.region
              awslogs-stream-prefix = "user"
              awslogs-create-group  = "true"
            }
          }

          environment = var.user_service_environment
        }
      }

      create_tasks_iam_role = false

      create_task_exec_iam_role = false
      task_exec_iam_role_arn    = "arn:aws:iam::281657013469:role/ecsTaskExecutionRole"
      service_connect_configuration = {
        namespace = module.ecs_execution_role.role_arn
        service = {
          client_alias = {
            port     = 8081
            dns_name = "user-service"
          }
          port_name      = "user-service"
          discovery_name = "user-service"
        }
      }
      load_balancer = {
        service = {
          target_group_arn = module.alb.target_groups
          container_name   = "user-service"
          container_port   = 8081
        }
      }
      subnet_ids            = module.vpc.vpc_private_subnets
      create_security_group = false
      security_group_ids    = [module.sg_ecs.sg_id]
    }
  }
}
