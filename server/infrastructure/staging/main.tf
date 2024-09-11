provider "aws" {
  region = var.region
}

// VPC module
module "vpc" {
  source = "./../modules/vpc"

  vpc_name        = var.vpc_name
  vpc_cidr        = var.vpc_cidr
  vpc_azs         = var.vpc_azs
  vpc_pub_subnets = var.vpc_pub_subnets
  vpc_pri_subnets = var.vpc_pri_subnets
}

// SG modules
module "sg_alb" {
  source = "./../modules/securitygroup"

  sg_vpc_id   = module.vpc.vpc_id
  sg_alb_name = var.sg_alb_name
}

module "sg_ecs" {
  source = "./../modules/securitygroup"

  sg_vpc_id   = module.vpc.vpc_id
  sg_ecs_name = var.sg_ecs_name
}

module "sg_rds" {
  source = "./../modules/securitygroup"

  sg_vpc_id   = module.vpc.vpc_id
  sg_rds_name = var.sg_rds_name
}

// Load Balancer module
module "alb" {
  source = "./../modules/loadbalancer"

  alb_name            = var.alb_name
  alb_subnets         = module.vpc.vpc_public_subnets
  alb_security_groups = [module.sg_alb.sg_alb_id] // List?
  alb_listener        = var.alb_listener
  alb_target_groups   = var.alb_target_groups
}
