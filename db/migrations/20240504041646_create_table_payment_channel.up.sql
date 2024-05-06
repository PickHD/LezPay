CREATE TYPE payment_channel_status AS ENUM ('ACTIVE','INACTIVE');

CREATE TABLE IF NOT EXISTS payment_channel (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR NOT NULL,
    code VARCHAR NOT NULL,
    image_url TEXT NULL,
    payment_instruction TEXT NULL,
    status payment_channel_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);