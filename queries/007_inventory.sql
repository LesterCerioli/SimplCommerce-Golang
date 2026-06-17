-- ============================================================
-- WAREHOUSES
-- ============================================================
CREATE TABLE IF NOT EXISTS inventory_warehouses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    address_id UUID REFERENCES identity_addresses(id) ON DELETE RESTRICT,
    vendor_id UUID REFERENCES identity_vendors(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE inventory_warehouses IS 'Warehouses for inventory management';

CREATE INDEX IF NOT EXISTS idx_inventory_warehouses_address_id ON inventory_warehouses(address_id);
CREATE INDEX IF NOT EXISTS idx_inventory_warehouses_vendor_id ON inventory_warehouses(vendor_id);

-- ============================================================
-- STOCKS
-- ============================================================
CREATE TABLE IF NOT EXISTS inventory_stocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    warehouse_id UUID NOT NULL REFERENCES inventory_warehouses(id) ON DELETE RESTRICT,
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved_quantity INT NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
    low_stock_threshold INT DEFAULT 5 CHECK (low_stock_threshold >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (product_id, warehouse_id)
);

COMMENT ON TABLE inventory_stocks IS 'Stock levels per product per warehouse';

CREATE INDEX IF NOT EXISTS idx_inventory_stocks_product_id ON inventory_stocks(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_stocks_warehouse_id ON inventory_stocks(warehouse_id);

-- ============================================================
-- STOCK HISTORIES
-- ============================================================
CREATE TABLE IF NOT EXISTS inventory_stock_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    warehouse_id UUID NOT NULL REFERENCES inventory_warehouses(id) ON DELETE RESTRICT,
    created_by_id UUID REFERENCES identity_users(id) ON DELETE RESTRICT,
    old_quantity INT,
    new_quantity INT NOT NULL,
    quantity_delta INT NOT NULL,
    reason VARCHAR(500),
    reference_type VARCHAR(100),
    reference_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE inventory_stock_histories IS 'Audit log of stock level changes';

CREATE INDEX IF NOT EXISTS idx_inventory_stock_histories_product_id ON inventory_stock_histories(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_stock_histories_warehouse_id ON inventory_stock_histories(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inventory_stock_histories_created_by_id ON inventory_stock_histories(created_by_id);

-- ============================================================
-- PRODUCT BACK IN STOCK SUBSCRIPTIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS inventory_product_back_in_stock_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    customer_email VARCHAR(255) NOT NULL,
    is_notified BOOLEAN DEFAULT false,
    notified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE inventory_product_back_in_stock_subscriptions IS 'Customer subscriptions for back-in-stock notifications';

CREATE INDEX IF NOT EXISTS idx_inventory_back_in_stock_product_id ON inventory_product_back_in_stock_subscriptions(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_back_in_stock_email ON inventory_product_back_in_stock_subscriptions(customer_email);

-- ============================================================
-- Add FK from shipping_shipments.warehouse_id (table in 006)
-- ============================================================
ALTER TABLE shipping_shipments ADD CONSTRAINT fk_shipments_warehouse
    FOREIGN KEY (warehouse_id) REFERENCES inventory_warehouses(id) ON DELETE RESTRICT;
