output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.vpc.vpc_id
}

output "vpc_cidr" {
  description = "The CIDR block of the VPC"
  value       = module.vpc.vpc_cidr_block
}

output "vpc_public_subnets" {
  description = "The IDs of the public subnets"
  value       = module.vpc.public_subnets
}

output "vpc_public_subnets_cidr_blocks" {
  value = module.vpc.public_subnets_cidr_blocks
}

output "vpc_private_subnets" {
  description = "The IDs of the private subnets"
  value       = module.vpc.private_subnets
}

output "vpc_private_subnets_cidr_blocks" {
  value = module.vpc.private_subnets_cidr_blocks
}

output "vpc_database_subnet_group" {
  description = "The vpc database subnet group"
  value       = module.vpc.database_subnet_group
}
