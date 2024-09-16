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

module "sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.2"

  vpc_id = var.sg_vpc_id

  name        = var.sg_name
  description = var.sg_description

  ingress_cidr_blocks                   = var.ingress_cidr_blocks
  ingress_rules                         = var.ingress_rules
  ingress_with_source_security_group_id = var.ingress_with_source_security_group_id

  egress_cidr_blocks = var.egress_cidr_blocks
  egress_rules       = var.egress_rules
}
