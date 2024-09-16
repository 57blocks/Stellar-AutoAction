variable "cluster_name" {
  description = "Name of the IAM role"
  type        = string
  default     = ""
}

variable "services" {
  description = "List of services to be created in the ECS cluster"
  type        = any
  default     = {}
}