# resource "aws_security_group" "this" {
#   name        = var.sg_name
#   vpc_id      = var.vpc_id
#   description = var.sg_description
#
#   tags = var.tags
# }
#
# resource "aws_security_group_rule" "ingress" {
#   type              = "ingress"
#   from_port         = var.ingress_from_port
#   to_port           = var.ingress_to_port
#   protocol          = var.ingress_protocol
#   cidr_blocks       = var.ingress_cidr_blocks
#   security_group_id = aws_security_group.this.id
# }
#
# resource "aws_security_group_rule" "egress" {
#   type              = "egress"
#   from_port         = var.egress_from_port
#   to_port           = var.egress_to_port
#   protocol          = var.egress_protocol
#   cidr_blocks       = var.egress_cidr_blocks
#   security_group_id = aws_security_group.this.id
# }

module "sg_alb" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.2"

  vpc_id = var.sg_vpc_id

  name        = var.sg_alb_name
  description = var.sg_alb_description

  ingress_cidr_blocks = var.ingress_cidr_blocks
  ingress_rules       = var.ingress_rules // why two 80s in the rule?

  egress_cidr_blocks = var.egress_cidr_blocks
  egress_rules       = var.egress_rules
}

module "sg_ecs" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.2"

  vpc_id = var.sg_vpc_id

  name        = var.sg_ecs_name
  description = "Security group for applications hosted on ECS"

  ingress_cidr_blocks                   = var.ingress_cidr_blocks
  ingress_with_source_security_group_id = [module.sg_alb.security_group_id]

  egress_cidr_blocks = var.egress_cidr_blocks
  egress_rules       = var.egress_rules
}

module "sg_rds" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.2"

  vpc_id = var.sg_vpc_id

  name        = var.sg_rds_name
  description = "Security group for PostgreSQL hosted on RDS"

  ingress_cidr_blocks = ["0.0.0.0/0"]
  ingress_rules       = ["postgresql-tcp"]

  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules       = ["all-all"]
}
