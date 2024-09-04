BEGIN;

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