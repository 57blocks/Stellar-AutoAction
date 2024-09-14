# module "iam_role" {
#   source  = "terraform-aws-modules/iam/aws"
#   version = "~>5.44.0"
#
#
#
resource "aws_iam_role" "this" {
  name               = var.role_name
  assume_role_policy = var.assume_role_policy
}

resource "aws_iam_role_policy" "this" {
  name   = var.role_policy_name
  role   = aws_iam_role.this.id
  policy = var.policy
}
