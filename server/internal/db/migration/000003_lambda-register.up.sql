BEGIN;

-- lambda info
DROP TABLE IF EXISTS "lambda";

CREATE TABLE "lambda" (
    "id" serial PRIMARY KEY,
    "function_name" varchar UNIQUE NOT NULL,
    "function_arn" varchar UNIQUE NOT NULL,
    "runtime" varchar NOT NULL,
    "timeout" int2 NOT NULL,
    "role" varchar NOT NULL,
    "handler" varchar NOT NULL,
    "description" varchar NOT NULL DEFAULT '',
    "code_sha256" varchar UNIQUE NOT NULL,
    "version" varchar NOT NULL,
    "revision_id" varchar UNIQUE NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    UNIQUE ("function_arn", "function_name")
);

CREATE INDEX ON "lambda" ("function_name");
CREATE INDEX ON "lambda" ("function_arn");

-- lambda scheduler info
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