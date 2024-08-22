BEGIN;

-- user
DROP TABLE IF EXISTS "principal_user";

CREATE TABLE "principal_user" (
    "id" serial PRIMARY KEY,
    "account" varchar UNIQUE NOT NULL,
    "password" text NOT NULL,
    "description" text NULL,
    "organization_id" integer NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "principal_user" ("account");

CREATE INDEX "principal_user_id_account_idx" ON "principal_user" ("id", "account");

-- organization
DROP TABLE IF EXISTS "principal_organization";

CREATE TABLE "principal_organization" (
    "id" serial PRIMARY KEY,
    "name" varchar UNIQUE NOT NULL,
    "description" text,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX "principal_organization_id_name_idx" ON "principal_organization" ("id", "name");

-- token
DROP TABLE IF EXISTS "principal_token";

CREATE TABLE "principal_token" (
    "id" serial PRIMARY KEY,
    "access" varchar UNIQUE NOT NULL,
    "refresh" varchar UNIQUE NOT NULL,
    "user_id" integer UNIQUE NOT NULL,
    "access_expires" timestamptz NOT NULL,
    "refresh_expires" timestamptz NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "principal_token" ("user_id");
CREATE INDEX ON "principal_token" ("access");
CREATE INDEX "principal_token_user_access_idx" ON "principal_token" ("user_id", "access");

COMMIT;