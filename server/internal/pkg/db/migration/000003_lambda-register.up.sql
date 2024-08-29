-- lambda info
DROP TABLE IF EXISTS "lambda";

CREATE TABLE "lambda" (
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

CREATE INDEX ON "lambda" ("function_arn");
CREATE INDEX ON "lambda" ("code_sha256");

CREATE INDEX "lambda_id_name_arn_idx" ON "lambda" ("id", "function_name", "function_arn");

-- vpc info which is the lambda function bound to
DROP TABLE IF EXISTS "lambda_vpc";

CREATE TABLE "lambda_vpc" (
    "id" serial PRIMARY KEY,
    "vpc_id" varchar UNIQUE NOT NULL,
    "subnet_ids" varchar[] NOT NULL,
    "security_group_ids" varchar[] NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "lambda_vpc" ("vpc_id");

