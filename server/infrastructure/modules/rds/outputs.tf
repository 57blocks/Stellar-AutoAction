output "db_endpoint" {
  value = module.rds.db_instance_endpoint
}

output "db_arn" {
  value = module.rds.db_instance_arn
}

output "db_address" {
  value = module.rds.db_instance_address
}
