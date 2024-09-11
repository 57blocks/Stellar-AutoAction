output "alb_dns" {
  description = "The DNS name of the ALB"
  value       = module.alb.dns_name
}