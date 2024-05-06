CREATE TABLE IF NOT EXISTS transaction_history (
    id SERIAL NOT NULL PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    status transaction_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);