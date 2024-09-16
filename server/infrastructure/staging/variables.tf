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

variable "alb_listener" {
  description = "ALB listener"
  type        = any
  default     = {}
}

variable "alb_target_groups" {
  description = "ALB target groups"
  type        = any
  default     = {}
}

variable "ecs_cluster_name" {
  description = "ECS cluster name"
  type        = string
  default     = ""
}
