BEGIN;

-- vpc info which is the lambda function bound to
DROP TABLE IF EXISTS "vpc";

CREATE TABLE "vpc" (
    "id" serial PRIMARY KEY,
    "organization_id" int4 NOT NULL,
    "aws_id" varchar UNIQUE NOT NULL,
    "subnet_ids" varchar[] NOT NULL,
    "security_group_ids" varchar[] NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "vpc" ("aws_id");
CREATE INDEX "vpc_org_id_aws_id_idx" ON "vpc" ("organization_id", "aws_id");

-- lambda info
DROP TABLE IF EXISTS "lambda";

CREATE TABLE "lambda" (
    "id" serial PRIMARY KEY,
    "vpc_id" int4 NOT NULL, -- vpc_id is the reference of the PK of vpc table instead of vpc_id in Amazon
    "function_name" varchar NOT NULL,
    "function_arn" varchar UNIQUE NOT NULL,
    "runtime" varchar NOT NULL,
    "role" varchar NOT NULL,
    "handler" varchar NOT NULL,
    "description" varchar NOT NULL DEFAULT '',
    "code_sha256" varchar UNIQUE NOT NULL,
    "version" varchar NOT NULL,
    "revision_id" varchar UNIQUE NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "lambda" ("function_name");
CREATE INDEX ON "lambda" ("function_arn");

CREATE INDEX "lambda_name_arn_idx" ON "lambda" ("function_name", "function_arn");
CREATE INDEX "lambda_name_version_idx" ON "lambda" ("function_name", "version");

-- vpc info which is the lambda function bound to
DROP TABLE IF EXISTS "lambda_scheduler";

CREATE TABLE "lambda_scheduler" (
    "id" serial PRIMARY KEY,
    "lambda_id" int4 NOT NULL,
    "schedule_arn" varchar UNIQUE NOT NULL,
    "expression" varchar NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "lambda_scheduler" ("lambda_id");

COMMIT;