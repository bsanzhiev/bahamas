CREATE TABLE IF NOT EXISTS account_holders (
    account_id UUID NOT NULL,
    user_id INT NOT NULL,
    PRIMARY KEY (account_id, user_id),
    FOREIGN KEY (account_id) REFERENCES accounts(account_id) ON DELETE CASCADE
);