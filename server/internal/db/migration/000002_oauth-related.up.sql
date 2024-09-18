BEGIN;

-- user
DROP TABLE IF EXISTS "user";

CREATE TABLE "user" (
    "id" serial PRIMARY KEY,
    "account" varchar NOT NULL,
    "password" text NOT NULL,
    "description" text NULL,
    "organization_id" integer NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    UNIQUE ("account", "organization_id")
);

CREATE INDEX ON "user" ("account");

-- organization
DROP TABLE IF EXISTS "organization";

CREATE TABLE "organization" (
    "id" serial PRIMARY KEY,
    "name" varchar UNIQUE NOT NULL,
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
    "access_id" varchar UNIQUE NOT NULL,
    "access_expires" timestamptz NOT NULL,
    "refresh" varchar UNIQUE NOT NULL,
    "refresh_id" varchar UNIQUE NOT NULL,
    "refresh_expires" timestamptz NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    UNIQUE ("user_id", "access")
);

CREATE INDEX ON "token" ("user_id");
CREATE INDEX ON "token" ("access_id");
CREATE INDEX ON "token" ("refresh_id");

COMMIT;