BEGIN;

-- user
DROP TABLE IF EXISTS "user";

CREATE TABLE "user" (
    "id" serial PRIMARY KEY,
    "account" varchar UNIQUE NOT NULL,
    "password" text NOT NULL,
    "description" text NULL,
    "organization_id" integer NOT NULL,
    "user_key" varchar NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "user" ("account");

CREATE INDEX "user_id_account_idx" ON "user" ("id", "account");

-- organization
DROP TABLE IF EXISTS "organization";

CREATE TABLE "organization" (
    "id" serial PRIMARY KEY,
    "name" varchar UNIQUE NOT NULL,
    "cube_signer_org" varchar UNIQUE NOT NULL,
    "description" text,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX "organization_id_name_idx" ON "organization" ("id", "name");

-- token
DROP TABLE IF EXISTS "token";

CREATE TABLE "token" (
    "id" serial PRIMARY KEY,
    "user_id" integer UNIQUE NOT NULL,
    "access" varchar UNIQUE NOT NULL,
    "refresh" varchar UNIQUE NOT NULL,
    "access_expires" timestamptz NOT NULL,
    "refresh_expires" timestamptz NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "token" ("user_id");
CREATE INDEX ON "token" ("access");
CREATE INDEX "token_user_access_idx" ON "token" ("user_id", "access");

COMMIT;