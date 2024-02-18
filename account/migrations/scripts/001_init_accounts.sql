CREATE TABLE IF NOT EXISTS accounts (
    account_id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    user_id INT NOT NULL,
    account_number VARCHAR(255) UNIQUE NOT NULL,
    account_type VARCHAR(255) NOT NULL,
    balance NUMERIC(15, 2) NOT NULL -- Предполагая точность до двух знаков после запятой и достаточно большой диапазон для баланса
);