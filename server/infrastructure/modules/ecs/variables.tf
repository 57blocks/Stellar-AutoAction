variable "ecs_cluster_name" {
  description = "Name of the IAM role"
  type        = string
  default     = ""
}

variable "ecs_fargate_capacity_providers" {
  description = "Map of Fargate capacity provider definitions to use for the cluster"
  type        = any
  default     = {}
}

variable "ecs_services" {
  description = "List of services to be created in the ECS cluster"
  type        = any
  default     = {}
}