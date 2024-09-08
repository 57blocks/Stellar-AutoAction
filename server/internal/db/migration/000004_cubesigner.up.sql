BEGIN;

-- CubeSigner Role information table
DROP TABLE IF EXISTS "cube_signer_role";

CREATE TABLE "cube_signer_role" (
    "id" serial PRIMARY KEY,
    "organization_id" int4 NOT NULL,
    "account_id" int4 NOT NULL,
    "role" varchar UNIQUE NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "cube_signer_role" ("organization_id");
CREATE INDEX ON "cube_signer_role" ("account_id");
CREATE INDEX "cube_signer_role_oar_idx" ON "cube_signer_role" ("organization_id", "account_id", "role");

-- CubeSigner Keys information table
DROP TABLE IF EXISTS "cube_signer_key";

CREATE TABLE "cube_signer_key" (
    "id" serial PRIMARY KEY,
    "role_id" varchar NOT NULL,
    "key" varchar UNIQUE NOT NULL,
    "scopes" varchar[] NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "cube_signer_key" ("role_id");
CREATE INDEX ON "cube_signer_key" ("key");
CREATE INDEX "cube_signer_key_rk_idx" ON "cube_signer_key" ("role_id", "key");

COMMIT;