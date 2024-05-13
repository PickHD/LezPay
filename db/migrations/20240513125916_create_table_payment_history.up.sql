CREATE TABLE IF NOT EXISTS payment_history (
    id SERIAL NOT NULL PRIMARY KEY,
    payment_id BIGINT NOT NULL,
    status payment_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);