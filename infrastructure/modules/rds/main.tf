module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 6.9.0"

  identifier = var.rds_identifier

  # All available versions: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_PostgreSQL.html#PostgreSQL.Concepts
  engine                   = "postgres"
  engine_version           = "16.4"
  engine_lifecycle_support = "open-source-rds-extended-support-disabled"
  family                   = "postgres16" # DB parameter group
  major_engine_version     = "16"         # DB option group
  instance_class           = "db.t3.micro"

  allocated_storage     = 20
  max_allocated_storage = 100

  # NOTE: Do NOT use 'user' as the value for 'username' as it throws:
  # "Error creating DB Instance: InvalidParameterValue: MasterUsername
  # user cannot be used as it is a reserved word used by the engine"
  db_name  = var.rds_db_name
  username = var.rds_username
  password = var.rds_password
  port     = 5432

  multi_az               = true
  vpc_security_group_ids = var.rds_vpc_security_group_ids

  db_subnet_group_name = var.rds_db_subnet_group_name
  subnet_ids           = var.rds_subnet_ids

  auto_minor_version_upgrade = true
  maintenance_window         = "Mon:00:00-Mon:03:00"

  create_cloudwatch_log_group     = true
  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]

  backup_window           = "03:00-05:00"
  backup_retention_period = 0

  // default: true. The password provided will not be used if `manage_master_user_password` is set to true.
  manage_master_user_password = false
  deletion_protection         = false
  skip_final_snapshot         = true
  publicly_accessible         = true

  parameters = [
    {
      name  = "autovacuum"
      value = 1
    },
    {
      name  = "client_encoding"
      value = "utf8"
    }
  ]
}