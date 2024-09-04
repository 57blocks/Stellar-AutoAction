BEGIN;

-- Organization secrets table
DROP TABLE IF EXISTS "organization_root_session";

CREATE TABLE "organization_root_session" (
    "id" serial PRIMARY KEY,
    "organization_id" int4 NOT NULL,
    "expiration" int4 NOT NULL,
    "token" text UNIQUE NOT NULL,
    "refresh_token" text UNIQUE NOT NULL,
    -- here blew is about the epoch log of the session
    "session_id" varchar UNIQUE NOT NULL,
    "epoch" int4 NOT NULL,
    "epoch_token" varchar UNIQUE NOT NULL,
    "epoch_auth_token" text UNIQUE NOT NULL,
    "epoch_auth_token_exp" int4 NOT NULL,
    "epoch_refresh_token" varchar UNIQUE NOT NULL,
    "epoch_refresh_token_exp" int4 NOT NULL,

    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "organization_root_session" ("organization_id");
CREATE INDEX ON "organization_root_session" ("session_id");
CREATE INDEX "organization_root_session_idx" ON
    "organization_root_session" ("organization_id", "session_id");

-- CubeSigner role-key information table
DROP TABLE IF EXISTS "organization_role_key";

CREATE TABLE "organization_role_key" (
    "id" serial PRIMARY KEY,
    "organization_id" int4 NOT NULL,
    "cs_role_id" varchar UNIQUE NOT NULL,
    "cs_key_id" varchar UNIQUE NOT NULL,
    "cs_scopes" varchar[] NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "organization_role_key" ("organization_id");
CREATE INDEX "organization_role_key_r_k_idx" ON "organization_role_key" ("organization_id", "cs_role_id", "cs_key_id");

COMMIT;