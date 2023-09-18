CREATE TYPE transaction_status AS ENUM ('PENDING','PAYMENT_PENDING','PAYMENT_EXPIRED','PAYMENT_FAILED','REFUND','SUCCESS');
CREATE TYPE transaction_type AS ENUM ('TOPUP','PAYOUT');

CREATE TABLE IF NOT EXISTS transaction (
    id BIGINT NOT NULL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    merchant_id BIGINT NOT NULL,
    payment_channel_id BIGINT NOT NULL,
    ref_id VARCHAR NULL,
    amount INT NOT NULL,
    status transaction_status NOT NULL,
    type transaction_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
)