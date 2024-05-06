CREATE TYPE transaction_status AS ENUM ('PENDING','PROCESS','EXPIRED','FAILED','REFUND','SUCCESS','FAILED_REFUND');
CREATE TYPE transaction_type AS ENUM ('TOPUP','PAYOUT', 'TRANSFER');

CREATE TABLE IF NOT EXISTS transaction (
    id SERIAL NOT NULL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    payment_channel_id BIGINT NOT NULL,
    merchant_id BIGINT NULL,
    ref_id VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    merchant_fee INT NULL,
    admin_fee INT NULL,
    products JSONB NULL,
    status transaction_status NOT NULL,
    type transaction_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);