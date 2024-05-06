CREATE TABLE IF NOT EXISTS wallet_usage_history (
  id SERIAL NOT NULL PRIMARY KEY,
  wallet_id BIGINT NOT NULL,
  balance_before DECIMAL(10,2) NOT NULL,
  balance_after DECIMAL(10,2) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NULL
);