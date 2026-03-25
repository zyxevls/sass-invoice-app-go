CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY,
    user_id UUID,
    invoice_code VARCHAR(50),
    client_email VARCHAR(100),
    status VARCHAR(20),
    total_amount BIGINT,
    expired_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS invoice_items (
    id UUID PRIMARY KEY,
    invoice_id UUID REFERENCES invoices(id) ON DELETE CASCADE,
    name VARCHAR(100),
    qty INT,
    price BIGINT
);