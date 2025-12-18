CREATE TABLE accounts (
    account_id SERIAL PRIMARY KEY,
    document_number VARCHAR(20) NOT NULL UNIQUE
);
