resource "aws_iam_role" "this" {
  name               = var.role_name
  description        = var.role_description
  assume_role_policy = var.assume_role_policy
}

resource "aws_iam_role_policy" "this" {
  name   = var.role_policy_name
  role   = aws_iam_role.this.id
  policy = var.policy
}
