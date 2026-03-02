CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    type VARCHAR(20) NOT NULL,
    amount BIGINT NOT NULL,
    category VARCHAR(100) NOT NULL,
    description TEXT,
    date TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);