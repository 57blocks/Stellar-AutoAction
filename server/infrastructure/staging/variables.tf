variable "region" {
  description = "The AWS region to deploy the VPC"
  type        = string
  default     = ""
}

variable "env" {
  description = "The environment where the VPC is being deployed"
  type        = string
  default     = ""
}

// VPC
variable "vpc_name" {
  description = "Vpc name"
  type        = string
  default     = ""
}

variable "vpc_cidr" {
  description = "The CIDR block for the VPC. Default value is a valid CIDR, but not acceptable by AWS and should be overriden"
  type        = string
  default     = ""
}

variable "vpc_azs" {
  description = "A list of availability zones in the region"
  type        = list(string)
  default     = []
}

variable "vpc_pub_subnets" {
  description = "A list of public subnets inside the VPC"
  type        = list(string)
  default     = []
}

variable "vpc_pri_subnets" {
  description = "A list of private subnets inside the VPC"
  type        = list(string)
  default     = []
}

// Security Groups
variable "sg_alb_name" {
  description = "ALB security group name"
  type        = string
  default     = ""
}

variable "sg_ecs_name" {
  description = "Application security group name"
  type        = string
  default     = ""
}

variable "sg_rds_name" {
  description = "RDS security group name"
  type        = string
  default     = ""
}

// Application Load Balancer
variable "alb_name" {
  description = "ALB name"
  type        = string
  default     = ""
}

# variable "alb_listener" {
#   description = "ALB listener"
#   type        = any
#   default     = {}
# }
#
# variable "alb_target_groups" {
#   description = "ALB target groups"
#   type        = any
#   default     = {}
# }

// Secrets Manager instances
variable "jwt_private_key" {
  description = "The private key name to sign JWT tokens"
  type        = string
  default     = ""
}

variable "jwt_private_key_value" {
  description = "The private key to sign JWT tokens"
  type        = string
  default     = ""
}

variable "jwt_public_key" {
  description = "The public key name to verify JWT tokens"
  type        = string
  default     = ""
}

variable "jwt_public_key_value" {
  description = "The public key to verify JWT tokens"
  type        = string
  default     = ""
}

variable "rds_password" {
  description = "The RDS password key name"
  type        = string
  default     = ""
}

variable "rds_password_value" {
  description = "The RDS password value"
  type        = string
  default     = ""
}

// RDS
variable "rds_identifier" {
  description = "The RDS instance identifier"
  type        = string
  default     = ""
}

variable "rds_db_name" {
  description = "The RDS database name"
  type        = string
  default     = ""
}

variable "rds_username" {
  description = "The RDS username"
  type        = string
  default     = ""
}

// ECR
variable "ecr_repository_name" {
  description = "The name of the repository"
  type        = string
  default     = ""
}

// ECS Cluster
variable "ecs_cluster_name" {
  description = "ECS cluster name"
  type        = string
  default     = ""
}

variable "ecs_fargate_capacity_providers" {
  description = "Map of Fargate capacity provider definitions to use for the cluster"
  type        = any
  default     = {}
}
