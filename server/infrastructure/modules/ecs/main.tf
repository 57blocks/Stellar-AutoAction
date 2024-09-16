# resource "aws_security_group" "services" {
#   name        = "services-security-group"
#   description = "All services traffic"
#   vpc_id      = var.vpc_id
#
#   ingress {
#     from_port       = 0
#     to_port         = 0
#     protocol        = "-1"
#     security_groups = [var.alb_security_group_id]
#   }
#
#   ingress {
#     from_port   = 0
#     to_port     = 0
#     protocol    = "-1"
#     cidr_blocks = [var.vpc_cidr_block]
#   }
#
#   egress {
#     from_port   = 0
#     to_port     = 0
#     protocol    = "-1"
#     cidr_blocks = ["0.0.0.0/0"]
#   }
# }
#
# resource "aws_service_discovery_http_namespace" "this" {
#   name        = var.cluster_name
#   description = "CloudMap namespace for ${var.cluster_name}"
# }

module "ecs" {
  source = "terraform-aws-modules/ecs/aws"

  cluster_name = var.cluster_name

  fargate_capacity_providers = {
    FARGATE = {
      default_capacity_provider_strategy = {
        weight = 50
        base   = 20
      }
    }
    FARGATE_SPOT = {
      default_capacity_provider_strategy = {
        weight = 50
      }
    }
  }

  services = var.services
}