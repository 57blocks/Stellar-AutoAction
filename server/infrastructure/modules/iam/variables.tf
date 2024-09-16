variable "role_name" {
  description = "Name of the IAM role"
  type        = string
  default     = ""
}

variable "role_description" {
  description = "Description of the IAM role"
  type        = string
  default     = ""
}

variable "assume_role_policy" {
  description = "Assume role policy document of the IAM role"
  type        = string
  default     = ""
}

variable "role_policy_name" {
  description = "Policy name of the IAM role"
  type        = string
  default     = ""
}

variable "policy" {
  description = "Policy of the IAM role"
  type        = string
  default     = ""
}
