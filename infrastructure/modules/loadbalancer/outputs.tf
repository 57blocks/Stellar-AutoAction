output "alb_dns" {
  description = "The DNS name of the ALB"
  value       = module.alb.dns_name
}

output "target_groups" {
  description = "The target groups of the ALB"
  value       = module.alb.target_groups["ecs_app"].arn
}
