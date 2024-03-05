CREATE TABLE IF NOT EXISTS transactions (
    transaction_id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    account_id UUID NOT NULL,
    recipient_account_id UUID NOT NULL,
    amount NUMERIC(12, 2) NOT NULL,
    -- Предполагая максимальную точность в 2 знака после запятой
    transaction_type VARCHAR(255) NOT NULL,
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,
    description TEXT,
    FOREIGN KEY (account_id) REFERENCES accounts(account_id),
    FOREIGN KEY (recipient_account_id) REFERENCES accounts(account_id)
);