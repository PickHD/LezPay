CREATE TABLE IF NOT EXISTS redeem_history (
    id SERIAL NOT NULL PRIMARY KEY,
    redeem_id BIGINT NOT NULL,
    status redeem_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);