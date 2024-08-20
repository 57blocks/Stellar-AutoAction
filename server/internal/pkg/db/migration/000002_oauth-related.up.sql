BEGIN;

-- user
DROP TABLE IF EXISTS "user";

CREATE TABLE "user" (
    "id" serial PRIMARY KEY,
    "account" varchar NOT NULL,
    "password" text NOT NULL,
    "description" text NULL,
    "organization_id" integer NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX ON "user" ("account");

CREATE INDEX "id_account_idx" ON "user" ("id", "account");

-- organization
DROP TABLE IF EXISTS "organization";

CREATE TABLE "organization" (
    "id" serial PRIMARY KEY,
    "name" varchar NOT NULL,
    "description" text,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX "id_name_idx" ON "organization" ("id", "name");

-- token
DROP TABLE IF EXISTS "token";

CREATE TABLE "token" (
    "access" varchar PRIMARY KEY,
    "refresh" varchar UNIQUE NOT NULL,
    "user_id" integer UNIQUE,
    "access_expires" timestamptz NOT NULL,
    "refresh_expires" timestamptz NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX ON "token" ("user_id");
CREATE INDEX "user_access_idx" ON "token" ("user_id", "access");

COMMIT;