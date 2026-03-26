CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    invoice_code VARCHAR(50) NOT NULL,
    client_email VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL,
    total_amount BIGINT NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS invoice_items (
    id UUID PRIMARY KEY NOT NULL,
    invoice_id UUID REFERENCES invoices(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    qty INT NOT NULL,
    price BIGINT NOT NULL
);