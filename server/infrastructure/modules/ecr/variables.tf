variable "ecr_repository_name" {
  description = "The name of the repository"
  type        = string
  default     = ""
}

variable "repository_read_write_access_arns" {
  description = "The ARNs of the IAM users/roles that have read/write access to the repository"
  type        = list(string)
  default     = []
}
