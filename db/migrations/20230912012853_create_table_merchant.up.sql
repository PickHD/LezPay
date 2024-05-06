CREATE TABLE IF NOT EXISTS merchant (
    id SERIAL NOT NULL PRIMARY KEY,
    full_name VARCHAR(50) NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    phone_number VARCHAR(30) NOT NULL,
    password VARCHAR(255) NOT NULL,
    callback_url TEXT NULL,
    bank_account VARCHAR(255) NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);