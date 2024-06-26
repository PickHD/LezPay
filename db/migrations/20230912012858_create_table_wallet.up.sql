CREATE TABLE IF NOT EXISTS wallet (
    id SERIAL NOT NULL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    balance DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
)