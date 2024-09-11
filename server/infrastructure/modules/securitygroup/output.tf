output "sg_alb_id" {
  description = "ID of the security group"
  value       = module.sg_alb.security_group_id
}

output "sg_ecs_id" {
  description = "ID of the security group"
  value       = module.sg_ecs.security_group_id
}

output "sg_rds_id" {
  description = "ID of the security group"
  value       = module.sg_rds.security_group_id
}