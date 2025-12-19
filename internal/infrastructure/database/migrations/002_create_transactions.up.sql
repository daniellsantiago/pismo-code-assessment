CREATE TABLE operation_types (
    operation_type_id INTEGER PRIMARY KEY,
    description VARCHAR(50) NOT NULL
);

INSERT INTO operation_types (operation_type_id, description) VALUES 
    (1, 'PURCHASE'),
    (2, 'INSTALLMENT PURCHASE'),
    (3, 'WITHDRAWAL'),
    (4, 'PAYMENT');

CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(account_id),
    operation_type_id INTEGER NOT NULL REFERENCES operation_types(operation_type_id),
    amount DECIMAL(15,2) NOT NULL,
    event_date TIMESTAMP NOT NULL DEFAULT NOW()
);
