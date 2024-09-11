module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.13"

  name = var.vpc_name
  azs  = var.vpc_azs
  cidr = var.vpc_cidr

  public_subnets  = var.vpc_pub_subnets
  private_subnets = var.vpc_pri_subnets

  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true
}

