-- +goose Up
CREATE TABLE orders (
    order_uuid VARCHAR(36) PRIMARY KEY,
    user_uuid VARCHAR(36) NOT NULL,
    part_uuids TEXT[], -- array для хранения UUID деталей
    total_price DECIMAL(10,2) NOT NULL,
    transaction_uuid VARCHAR(36), 
    payment_method VARCHAR(20) NOT NULL CHECK (payment_method IN ('UNKNOWN', 'CARD', 'SBP', 'CREDIT_CARD', 'INVESTOR_MONEY')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELED')),
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP 
);

CREATE INDEX idx_orders_user_uuid ON orders(user_uuid);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_payment_method ON orders(payment_method);

-- +goose Down
DROP TABLE orders;