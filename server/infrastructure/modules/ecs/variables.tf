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

# variable "ecs_services" {
#   description = "List of services to be created in the ECS cluster"
#   type        = any
#   default     = {}
# }

variable "ecs_services" {
  description = "Definition of ECS services and associated configurations"
  type = map(object({
    cpu    = number
    memory = number

    container_definitions = map(object({
      cpu                = number
      memory             = number
      essential          = bool
      image              = string
      memory_reservation = number
      port_mappings = list(object({
        containerPort = number
        hostPort      = number
        protocol      = string
        appProtocol   = string
      }))
      readonly_root_filesystem = bool
      log_configuration = object({
        logDriver = string
        options   = map(string)
      })
      environment = list(object({
        name  = string
        value = string
      }))
      secrets = list(object({
        name      = string
        valueFrom = string
      }))
    }))

    create_tasks_iam_role     = bool
    create_task_exec_iam_role = bool
    task_exec_iam_role_arn    = string

    subnet_ids = list(string)

    load_balancer = object({
      service = object({
        target_group_arn = string
        container_name   = string
        container_port   = number
      })
    })

    create_security_group = bool
    security_group_ids    = list(string)
  }))
}
