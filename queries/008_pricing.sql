-- ============================================================
-- CART RULES
-- ============================================================
CREATE TABLE IF NOT EXISTS pricing_cart_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    priority INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    is_coupon_required BOOLEAN DEFAULT false,
    usage_limit_per_coupon INT CHECK (usage_limit_per_coupon >= 0),
    usage_limit_per_user INT CHECK (usage_limit_per_user >= 0),
    discount_type VARCHAR(50) NOT NULL CHECK (discount_type IN ('percentage', 'fixed_amount', 'free_shipping')),
    discount_value DECIMAL(18,2) CHECK (discount_value >= 0),
    min_order_amount DECIMAL(18,2) CHECK (min_order_amount >= 0),
    max_discount_amount DECIMAL(18,2) CHECK (max_discount_amount >= 0),
    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE pricing_cart_rules IS 'Cart-level pricing rules and promotions';

-- ============================================================
-- COUPONS
-- ============================================================
CREATE TABLE IF NOT EXISTS pricing_coupons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_rule_id UUID NOT NULL REFERENCES pricing_cart_rules(id) ON DELETE RESTRICT,
    code VARCHAR(255) NOT NULL UNIQUE,
    usage_count INT DEFAULT 0 CHECK (usage_count >= 0),
    max_usage INT CHECK (max_usage >= 0),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE pricing_coupons IS 'Coupon codes linked to cart rules';

CREATE INDEX IF NOT EXISTS idx_pricing_coupons_cart_rule_id ON pricing_coupons(cart_rule_id);
CREATE INDEX IF NOT EXISTS idx_pricing_coupons_code ON pricing_coupons(code);

-- ============================================================
-- CART RULE CATEGORIES (composite PK)
-- ============================================================
CREATE TABLE IF NOT EXISTS pricing_cart_rule_categories (
    cart_rule_id UUID NOT NULL REFERENCES pricing_cart_rules(id) ON DELETE RESTRICT,
    category_id UUID NOT NULL REFERENCES catalog_categories(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (cart_rule_id, category_id)
);

COMMENT ON TABLE pricing_cart_rule_categories IS 'Categories that a cart rule applies to';

CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_categories_category_id ON pricing_cart_rule_categories(category_id);

-- ============================================================
-- CART RULE PRODUCTS (composite PK)
-- ============================================================
CREATE TABLE IF NOT EXISTS pricing_cart_rule_products (
    cart_rule_id UUID NOT NULL REFERENCES pricing_cart_rules(id) ON DELETE RESTRICT,
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (cart_rule_id, product_id)
);

COMMENT ON TABLE pricing_cart_rule_products IS 'Products that a cart rule applies to';

CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_products_product_id ON pricing_cart_rule_products(product_id);

-- ============================================================
-- CART RULE CUSTOMER GROUPS (composite PK)
-- ============================================================
CREATE TABLE IF NOT EXISTS pricing_cart_rule_customer_groups (
    cart_rule_id UUID NOT NULL REFERENCES pricing_cart_rules(id) ON DELETE RESTRICT,
    customer_group_id UUID NOT NULL REFERENCES identity_customer_groups(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (cart_rule_id, customer_group_id)
);

COMMENT ON TABLE pricing_cart_rule_customer_groups IS 'Customer groups that a cart rule applies to';

CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_cust_group_customer_group_id ON pricing_cart_rule_customer_groups(customer_group_id);

-- ============================================================
-- CART RULE USAGES
-- ============================================================
CREATE TABLE IF NOT EXISTS pricing_cart_rule_usages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_rule_id UUID NOT NULL REFERENCES pricing_cart_rules(id) ON DELETE RESTRICT,
    coupon_id UUID REFERENCES pricing_coupons(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
    order_id UUID NOT NULL REFERENCES orders_orders(id) ON DELETE RESTRICT,
    discount_amount DECIMAL(18,2) NOT NULL CHECK (discount_amount >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE pricing_cart_rule_usages IS 'Record of cart rule/coupon usage per order';

CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_usages_cart_rule_id ON pricing_cart_rule_usages(cart_rule_id);
CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_usages_coupon_id ON pricing_cart_rule_usages(coupon_id);
CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_usages_user_id ON pricing_cart_rule_usages(user_id);
CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_usages_order_id ON pricing_cart_rule_usages(order_id);

-- ============================================================
-- CATALOG RULES
-- ============================================================
CREATE TABLE IF NOT EXISTS pricing_catalog_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    priority INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    discount_type VARCHAR(50) NOT NULL CHECK (discount_type IN ('percentage', 'fixed_amount')),
    discount_value DECIMAL(18,2) NOT NULL CHECK (discount_value >= 0),
    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE pricing_catalog_rules IS 'Catalog-level pricing rules';

-- ============================================================
-- CATALOG RULE CUSTOMER GROUPS (composite PK)
-- ============================================================
CREATE TABLE IF NOT EXISTS pricing_catalog_rule_customer_groups (
    catalog_rule_id UUID NOT NULL REFERENCES pricing_catalog_rules(id) ON DELETE RESTRICT,
    customer_group_id UUID NOT NULL REFERENCES identity_customer_groups(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (catalog_rule_id, customer_group_id)
);

COMMENT ON TABLE pricing_catalog_rule_customer_groups IS 'Customer groups that a catalog rule applies to';

CREATE INDEX IF NOT EXISTS idx_pricing_catalog_rule_cust_group_customer_group_id ON pricing_catalog_rule_customer_groups(customer_group_id);
