module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 9.11.0"

  name                       = var.alb_name
  load_balancer_type         = "application"
  vpc_id                     = var.alb_vpc_id
  subnets                    = var.alb_subnets
  security_groups            = var.alb_security_groups
  listeners                  = var.alb_listener
  target_groups              = var.alb_target_groups
  enable_deletion_protection = false
}
