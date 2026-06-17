-- ============================================================
-- CART ITEMS
-- ============================================================
CREATE TABLE IF NOT EXISTS cart_cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    customer_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
    vendor_id UUID REFERENCES identity_vendors(id) ON DELETE RESTRICT,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    price DECIMAL(18,2) NOT NULL CHECK (price >= 0),
    custom_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE cart_cart_items IS 'Shopping cart items for customers';

CREATE INDEX IF NOT EXISTS idx_cart_cart_items_product_id ON cart_cart_items(product_id);
CREATE INDEX IF NOT EXISTS idx_cart_cart_items_customer_id ON cart_cart_items(customer_id);
CREATE INDEX IF NOT EXISTS idx_cart_cart_items_vendor_id ON cart_cart_items(vendor_id);
