-- migrations/000001_create_secrets_table.up.sql

CREATE TABLE secrets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    kv_data JSONB NOT NULL
);
