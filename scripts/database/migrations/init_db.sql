CREATE TABLE accounts
(
    id              VARCHAR(50) NOT NULL,
    document_number VARCHAR(50),
    created_at      TIMESTAMP   NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE transactions
(
    id                INT GENERATED ALWAYS AS IDENTITY,
    account_id        VARCHAR(50),
    operation_type_id INT NOT NULL,
    amount            FLOAT,
    event_date        TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_transaction_account
        FOREIGN KEY (account_id)
            REFERENCES accounts (id)
            ON DELETE CASCADE
);