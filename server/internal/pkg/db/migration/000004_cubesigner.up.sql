-- cubesigner information table
DROP TABLE IF EXISTS "principal_org_key_pairs";

CREATE TABLE "principal_org_key_pairs" (
    "id" serial PRIMARY KEY,
    "organization_id" varchar NOT NULL,
    "public_key" varchar UNIQUE NOT NULL,
    "private_key" varchar UNIQUE NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP(2) NOT NULL
);

CREATE INDEX ON "principal_org_key_pairs" ("organization_id");
CREATE INDEX "principal_org_key_pairs_pub_pri_idx" ON "principal_org_key_pairs" ("public_key", "private_key");