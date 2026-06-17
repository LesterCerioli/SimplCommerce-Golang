-- ============================================================
-- PAYMENT PROVIDERS
-- ============================================================
CREATE TABLE IF NOT EXISTS payments_payment_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    config JSONB,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE payments_payment_providers IS 'Configured payment providers/gateways';

-- ============================================================
-- PAYMENTS
-- ============================================================
CREATE TABLE IF NOT EXISTS payments_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders_orders(id) ON DELETE RESTRICT,
    payment_provider_id UUID REFERENCES payments_payment_providers(id) ON DELETE RESTRICT,
    transaction_id VARCHAR(500),
    amount DECIMAL(18,2) NOT NULL CHECK (amount >= 0),
    currency_code VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    payment_method VARCHAR(255),
    gateway_response JSONB,
    is_refunded BOOLEAN DEFAULT false,
    refunded_amount DECIMAL(18,2) DEFAULT 0 CHECK (refunded_amount >= 0),
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE payments_payments IS 'Payment transactions for orders';

CREATE INDEX IF NOT EXISTS idx_payments_payments_order_id ON payments_payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_payments_payment_provider_id ON payments_payments(payment_provider_id);
CREATE INDEX IF NOT EXISTS idx_payments_payments_transaction_id ON payments_payments(transaction_id);
