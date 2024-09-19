module "sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.2.0"

  vpc_id = var.sg_vpc_id

  name        = var.sg_name
  description = var.sg_description

  ingress_cidr_blocks                   = var.ingress_cidr_blocks
  ingress_rules                         = var.ingress_rules
  ingress_with_source_security_group_id = var.ingress_with_source_security_group_id

  egress_cidr_blocks = var.egress_cidr_blocks
  egress_rules       = var.egress_rules
}
