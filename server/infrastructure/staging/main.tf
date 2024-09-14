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

  sg_vpc_id   = module.vpc.vpc_id
  sg_alb_name = var.sg_alb_name
  sg_alb_description = "Security group for ALB"

  ingress_cidr_blocks = ["0.0.0.0/0"]
  ingress_rules       = ["http-80-tcp", "https-443-tcp"] // why two 80s in the rule?

  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules       = ["all-all"]
}

module "sg_ecs" {
  source = "./../modules/securitygroup"

  sg_vpc_id   = module.vpc.vpc_id
  sg_ecs_name = var.sg_ecs_name

  ingress_with_source_security_group_id = [
    {
#       description              = "http from service two"
#       rule                     = "http-80-tcp"
      source_security_group_id = module.sg_alb.sg_alb_id
    }
  ]

  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules       = ["all-all"]
}

module "sg_rds" {
  source = "./../modules/securitygroup"

  // TODO
  sg_vpc_id   = module.vpc.vpc_id
  sg_rds_name = var.sg_rds_name
}

// Load Balancer module
module "alb" {
  source = "./../modules/loadbalancer"

  alb_name            = var.alb_name
  alb_subnets         = module.vpc.vpc_public_subnets
  alb_security_groups = [module.sg_alb.sg_alb_id] // List?
  alb_listener        = var.alb_listener
  alb_target_groups   = var.alb_target_groups
}

// IAM roles
module "lambda_execution_role" {
  source = "../modules/iam"

  role_name = "lambda_execution_role"
  assume_role_policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "events.amazonaws.com",
            "lambda.amazonaws.com"
          ]
        },
        "Action": "sts:AssumeRole"
      }
    ]
  })

  role_policy_name = "lambda_execution_role_policy"
  policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "logs:*",
          "s3:GetObject",
          "s3:PutObject",
          "secretsmanager:*"
        ],
        "Resource": "arn:aws:*:*:*:*"
      }
    ]
  })
}

module "ecs_execution_role" {
  source = "../modules/iam"

  role_name = "ecs_execution_role"
  assume_role_policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "",
        "Effect": "Allow",
        "Principal": {
          "Service": "ecs-tasks.amazonaws.com"
        },
        "Action": "sts:AssumeRole"
      }
    ]
  })

  role_policy_name = "lambda_execution_role_policy"
  policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          // TODO, minimal the scopes of ecr
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:GetRepositoryPolicy",
          "ecr:DescribeRepositories",
          "ecr:ListImages",
          "ecr:DescribeImages",
          "ecr:BatchGetImage",
          "ecr:GetLifecyclePolicy",
          "ecr:GetLifecyclePolicyPreview",
          "ecr:ListTagsForResource",
          "ecr:DescribeImageScanFindings",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload",
          "ecr:PutImage",
          "logs:*",
          "cloudwatch:GenerateQuery"
        ],
        "Resource": "*"
      }
    ]
  })
}

module "scheduler_invoke_role" {
  source = "../modules/iam"

  role_name = "scheduler_invoke_role"
}
