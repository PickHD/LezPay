CREATE TYPE redeem_status AS ENUM ('PENDING','PROCESS','SUCCESS','FAILED');

CREATE TABLE IF NOT EXISTS redeem (
    id SERIAL NOT NULL PRIMARY KEY,
    merchant_id BIGINT NOT NULL,
    amount INT NOT NULL,
    status redeem_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);