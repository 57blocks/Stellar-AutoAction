module "iam" {
  source  = "terraform-aws-modules/iam/aws"
  version = "~>5.40.0"

  name = var.iam_name
}