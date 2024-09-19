variable "rds_identifier" {
  description = "RDS instance identifier"
  type        = string
  default     = ""
}

variable "rds_db_name" {
  description = "RDS database name"
  type        = string
  default     = ""
}

variable "rds_username" {
  description = "RDS username"
  type        = string
  default     = ""
}

variable "rds_password" {
  description = "RDS password"
  type        = string
  default     = ""
}

variable "rds_db_subnet_group" {
  description = "RDS database subnet group in VPC"
  type        = string
  default     = ""
}

variable "rds_vpc_security_group_ids" {
  description = "RDS VPC security group IDs"
  type        = list(string)
  default     = []
}

variable "rds_subnet_ids" {
  description = "A list of VPC subnet IDs"
  type        = list(string)
  default     = []
}
