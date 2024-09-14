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

variable "sg_alb_description" {
  description = "ALB security group description"
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

variable "ingress_cidr_blocks" {
    description = "CIDR blocks to allow ingress traffic"
    type        = list(string)
    default     = []
}

variable "ingress_rules" {
    description = "List of ingress rules"
    type        = list(string)
    default     = []
}

variable "ingress_with_source_security_group_id" {
  description = "List of ingress rules with source security group id"
  type        = list(map(string))
  default     = []
}

variable "egress_cidr_blocks" {
    description = "CIDR blocks to allow egress traffic"
    type        = list(string)
    default     = []
}

variable "egress_rules" {
    description = "List of egress rules"
    type        = list(string)
    default     = []
}
