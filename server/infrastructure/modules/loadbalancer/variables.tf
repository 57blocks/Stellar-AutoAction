# variable "region" {
#   description = "The AWS region to deploy the VPC"
#   type        = string
#   default     = ""
# }

variable "alb_name" {
  description = "Name to be used on all the resources as identifier"
  type        = string
  default     = ""
}

variable "alb_vpc_id" {
  description = "VPC in for security group"
  type        = string
  default     = ""
}

variable "alb_subnets" {
  description = "The security group for the ALB"
  type        = list(string)
  default     = []
}

variable "alb_security_groups" {
  description = "The security group for the ALB"
  type        = list(string)
  default     = []
}

variable "alb_listener" {
  description = "The ALB listener configuration"
  type        = any
  default     = {}
}

variable "alb_target_groups" {
  description = "The target group configuration"
  type        = any
  default     = {}
}
