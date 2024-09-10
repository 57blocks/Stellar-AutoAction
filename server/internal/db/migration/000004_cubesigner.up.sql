BEGIN;

-- CubeSigner Keys information table
DROP TABLE IF EXISTS "cube_signer_key";

CREATE TABLE "cube_signer_key" (
    "id" serial PRIMARY KEY,
    "account_id" int4 NOT NULL,
    "key" varchar UNIQUE NOT NULL,
    "scopes" varchar[] NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    UNIQUE ("account_id", "key")
);

CREATE INDEX ON "cube_signer_key" ("account_id");
CREATE INDEX ON "cube_signer_key" ("key");

COMMIT;