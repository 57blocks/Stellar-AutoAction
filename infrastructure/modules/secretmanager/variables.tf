variable "secret_name" {
  description = "Name of the secret."
  default     = "example-secret"
}

variable "secret_value" {
  description = "The value of the secret."
  sensitive   = true
  default     = "my-secret-value"
}
