-- lambda info
DROP TABLE IF EXISTS "object_lambda";

CREATE TABLE "object_lambda" (
    "id" serial PRIMARY KEY,
    "function_name" varchar NOT NULL,
    "function_arn" varchar UNIQUE NOT NULL,
    "runtime" varchar NOT NULL,
    "role" varchar NOT NULL,
    "handler" varchar NOT NULL,
    "description" varchar NOT NULL DEFAULT '',
    "code_sha256" varchar UNIQUE NOT NULL,
    "version" varchar NOT NULL,
    "revision_id" varchar UNIQUE NOT NULL,
    "vpc_bound" boolean NOT NULL DEFAULT FALSE,
    "vpc_id" varchar UNIQUE NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "subject_lambda" ("function_arn");
CREATE INDEX ON "subject_lambda" ("code_sha256");

CREATE INDEX "subject_lambda_id_name_arn_idx" ON "object_lambda" ("id", "function_name", "function_arn");

-- vpc info which is the lambda function bound to
DROP TABLE IF EXISTS "subject_lambda_vpc";

CREATE TABLE "object_lambda_vpc" (
    "id" serial PRIMARY KEY,
    "vpc_id" varchar UNIQUE NOT NULL,
    "subnet_ids" varchar[] NOT NULL,
    "security_group_ids" varchar[] NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "object_lambda_vpc" ("vpc_id");

