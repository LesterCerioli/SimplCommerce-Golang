-- ============================================================
-- SHIPPING PROVIDERS
-- ============================================================
CREATE TABLE IF NOT EXISTS shipping_shipping_providers (
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

COMMENT ON TABLE shipping_shipping_providers IS 'Configured shipping providers/carriers';

-- ============================================================
-- SHIPMENTS
-- ============================================================
CREATE TABLE IF NOT EXISTS shipping_shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders_orders(id) ON DELETE RESTRICT,
    warehouse_id UUID,
    vendor_id UUID REFERENCES identity_vendors(id) ON DELETE RESTRICT,
    created_by_id UUID REFERENCES identity_users(id) ON DELETE RESTRICT,
    shipping_provider_id UUID REFERENCES shipping_shipping_providers(id) ON DELETE RESTRICT,
    tracking_number VARCHAR(500),
    tracking_url TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    shipping_method VARCHAR(255),
    weight DECIMAL(10,2) CHECK (weight >= 0),
    shipping_cost DECIMAL(18,2) CHECK (shipping_cost >= 0),
    shipped_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE shipping_shipments IS 'Shipments for order fulfillment';

CREATE INDEX IF NOT EXISTS idx_shipping_shipments_order_id ON shipping_shipments(order_id);
CREATE INDEX IF NOT EXISTS idx_shipping_shipments_warehouse_id ON shipping_shipments(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_shipping_shipments_vendor_id ON shipping_shipments(vendor_id);
CREATE INDEX IF NOT EXISTS idx_shipping_shipments_created_by_id ON shipping_shipments(created_by_id);
CREATE INDEX IF NOT EXISTS idx_shipping_shipments_shipping_provider_id ON shipping_shipments(shipping_provider_id);
CREATE INDEX IF NOT EXISTS idx_shipping_shipments_tracking_number ON shipping_shipments(tracking_number);

-- ============================================================
-- SHIPMENT ITEMS
-- ============================================================
CREATE TABLE IF NOT EXISTS shipping_shipment_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_id UUID NOT NULL REFERENCES shipping_shipments(id) ON DELETE RESTRICT,
    order_item_id UUID NOT NULL REFERENCES orders_order_items(id) ON DELETE RESTRICT,
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE shipping_shipment_items IS 'Items included in a shipment';

CREATE INDEX IF NOT EXISTS idx_shipping_shipment_items_shipment_id ON shipping_shipment_items(shipment_id);
CREATE INDEX IF NOT EXISTS idx_shipping_shipment_items_order_item_id ON shipping_shipment_items(order_item_id);
CREATE INDEX IF NOT EXISTS idx_shipping_shipment_items_product_id ON shipping_shipment_items(product_id);

-- ============================================================
-- TABLE RATE PRICE AND DESTINATIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS shipping_table_rate_price_and_destinations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipping_provider_id UUID REFERENCES shipping_shipping_providers(id) ON DELETE RESTRICT,
    country_id UUID REFERENCES identity_countries(id) ON DELETE RESTRICT,
    state_or_province_id UUID REFERENCES identity_state_or_provinces(id) ON DELETE RESTRICT,
    district_id UUID REFERENCES identity_districts(id) ON DELETE RESTRICT,
    min_weight DECIMAL(10,2) CHECK (min_weight >= 0),
    max_weight DECIMAL(10,2) CHECK (max_weight >= 0),
    min_price DECIMAL(18,2) CHECK (min_price >= 0),
    max_price DECIMAL(18,2) CHECK (max_price >= 0),
    shipping_cost DECIMAL(18,2) NOT NULL CHECK (shipping_cost >= 0),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE shipping_table_rate_price_and_destinations IS 'Shipping rate rules based on destination and price';

CREATE INDEX IF NOT EXISTS idx_shipping_table_rate_shipping_provider_id ON shipping_table_rate_price_and_destinations(shipping_provider_id);
CREATE INDEX IF NOT EXISTS idx_shipping_table_rate_country_id ON shipping_table_rate_price_and_destinations(country_id);
CREATE INDEX IF NOT EXISTS idx_shipping_table_rate_state_or_province_id ON shipping_table_rate_price_and_destinations(state_or_province_id);
CREATE INDEX IF NOT EXISTS idx_shipping_table_rate_district_id ON shipping_table_rate_price_and_destinations(district_id);
