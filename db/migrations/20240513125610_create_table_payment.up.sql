CREATE TYPE payment_status AS ENUM ('PENDING','PROCESS','EXPIRED','FAILED','REFUND','SUCCESS','FAILED_REFUND');

CREATE TABLE IF NOT EXISTS payment (
    id SERIAL NOT NULL PRIMARY KEY,
    payment_channel VARCHAR(255) NOT NULL,
    ref_id VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    merchant_fee INT NULL,
    admin_fee INT NULL,
    payment_details JSONB NULL,
    status payment_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);