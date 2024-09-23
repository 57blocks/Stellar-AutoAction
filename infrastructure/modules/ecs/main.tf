module "ecs" {
  source  = "terraform-aws-modules/ecs/aws"
  version = "~> 5.11.0"

  cluster_name = var.ecs_cluster_name

  fargate_capacity_providers = var.ecs_fargate_capacity_providers

  services = var.ecs_services
}