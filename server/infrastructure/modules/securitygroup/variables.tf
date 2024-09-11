variable "sg_vpc_id" {
  description = "VPC in for security group"
  type        = string
  default     = ""
}

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
