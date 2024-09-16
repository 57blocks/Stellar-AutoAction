output "sg_id" {
  description = "ID of the security group"
  value       = module.sg.security_group_id
}

output "sg_arn" {
  description = "ARN of the security group"
  value       = module.sg.security_group_arn
}
