output "secret_arn" {
  value = aws_secretsmanager_secret.secret.arn
}

output "secret_value_id" {
  value = aws_secretsmanager_secret_version.secret_version.id
}

output "secret_value" {
  value = aws_secretsmanager_secret_version.secret_version.secret_string
}
