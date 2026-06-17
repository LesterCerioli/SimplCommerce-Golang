-- ============================================================
-- ORDERS
-- ============================================================
CREATE TABLE IF NOT EXISTS orders_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
    shipping_address_id UUID REFERENCES identity_addresses(id) ON DELETE RESTRICT,
    billing_address_id UUID REFERENCES identity_addresses(id) ON DELETE RESTRICT,
    parent_id UUID REFERENCES orders_orders(id) ON DELETE RESTRICT,
    order_number VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    currency_code VARCHAR(3) NOT NULL DEFAULT 'USD',
    subtotal DECIMAL(18,2) NOT NULL CHECK (subtotal >= 0),
    discount_total DECIMAL(18,2) DEFAULT 0 CHECK (discount_total >= 0),
    tax_total DECIMAL(18,2) DEFAULT 0 CHECK (tax_total >= 0),
    shipping_total DECIMAL(18,2) DEFAULT 0 CHECK (shipping_total >= 0),
    grand_total DECIMAL(18,2) NOT NULL CHECK (grand_total >= 0),
    paid_total DECIMAL(18,2) DEFAULT 0 CHECK (paid_total >= 0),
    refunded_total DECIMAL(18,2) DEFAULT 0 CHECK (refunded_total >= 0),
    customer_notes TEXT,
    staff_notes TEXT,
    shipping_method VARCHAR(255),
    payment_method VARCHAR(255),
    is_gift BOOLEAN DEFAULT false,
    gift_message TEXT,
    ordered_at TIMESTAMPTZ,
    shipped_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE orders_orders IS 'Customer orders';

CREATE INDEX IF NOT EXISTS idx_orders_orders_customer_id ON orders_orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_orders_shipping_address_id ON orders_orders(shipping_address_id);
CREATE INDEX IF NOT EXISTS idx_orders_orders_billing_address_id ON orders_orders(billing_address_id);
CREATE INDEX IF NOT EXISTS idx_orders_orders_parent_id ON orders_orders(parent_id);
CREATE INDEX IF NOT EXISTS idx_orders_orders_order_number ON orders_orders(order_number);
CREATE INDEX IF NOT EXISTS idx_orders_orders_status ON orders_orders(status);

-- ============================================================
-- ORDER ITEMS
-- ============================================================
CREATE TABLE IF NOT EXISTS orders_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders_orders(id) ON DELETE RESTRICT,
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    product_name VARCHAR(500) NOT NULL,
    product_sku VARCHAR(255) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(18,2) NOT NULL CHECK (unit_price >= 0),
    discount_total DECIMAL(18,2) DEFAULT 0 CHECK (discount_total >= 0),
    tax_total DECIMAL(18,2) DEFAULT 0 CHECK (tax_total >= 0),
    total DECIMAL(18,2) NOT NULL CHECK (total >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE orders_order_items IS 'Line items within an order';

CREATE INDEX IF NOT EXISTS idx_orders_order_items_order_id ON orders_order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_orders_order_items_product_id ON orders_order_items(product_id);

-- ============================================================
-- ORDER ADDRESSES
-- ============================================================
CREATE TABLE IF NOT EXISTS orders_order_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders_orders(id) ON DELETE RESTRICT,
    address_type VARCHAR(50) NOT NULL CHECK (address_type IN ('shipping', 'billing')),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    company VARCHAR(255),
    address_line1 VARCHAR(500) NOT NULL,
    address_line2 VARCHAR(500),
    city VARCHAR(255) NOT NULL,
    state VARCHAR(255),
    zip_code VARCHAR(20),
    country VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    email VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE orders_order_addresses IS 'Snapshots of addresses at time of order';

CREATE INDEX IF NOT EXISTS idx_orders_order_addresses_order_id ON orders_order_addresses(order_id);

-- ============================================================
-- ORDER HISTORIES
-- ============================================================
CREATE TABLE IF NOT EXISTS orders_order_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders_orders(id) ON DELETE RESTRICT,
    old_status VARCHAR(50),
    new_status VARCHAR(50) NOT NULL,
    notes TEXT,
    created_by_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE orders_order_histories IS 'Status change history for orders';

CREATE INDEX IF NOT EXISTS idx_orders_order_histories_order_id ON orders_order_histories(order_id);
CREATE INDEX IF NOT EXISTS idx_orders_order_histories_created_by_id ON orders_order_histories(created_by_id);
