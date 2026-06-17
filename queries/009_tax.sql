-- ============================================================
-- TAX CLASSES
-- ============================================================
CREATE TABLE IF NOT EXISTS tax_tax_classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE tax_tax_classes IS 'Tax classes for product categorization';

-- ============================================================
-- TAX RATES
-- ============================================================
CREATE TABLE IF NOT EXISTS tax_tax_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tax_class_id UUID NOT NULL REFERENCES tax_tax_classes(id) ON DELETE RESTRICT,
    country_id UUID NOT NULL REFERENCES identity_countries(id) ON DELETE RESTRICT,
    state_or_province_id UUID REFERENCES identity_state_or_provinces(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    rate DECIMAL(5,4) NOT NULL CHECK (rate >= 0),
    is_compound BOOLEAN DEFAULT false,
    priority INT DEFAULT 0,
    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE tax_tax_rates IS 'Tax rates per class per location';

CREATE INDEX IF NOT EXISTS idx_tax_tax_rates_tax_class_id ON tax_tax_rates(tax_class_id);
CREATE INDEX IF NOT EXISTS idx_tax_tax_rates_country_id ON tax_tax_rates(country_id);
CREATE INDEX IF NOT EXISTS idx_tax_tax_rates_state_or_province_id ON tax_tax_rates(state_or_province_id);

-- ============================================================
-- Add FK from catalog_products.tax_class_id (table in 002)
-- ============================================================
ALTER TABLE catalog_products ADD CONSTRAINT fk_products_tax_class
    FOREIGN KEY (tax_class_id) REFERENCES tax_tax_classes(id) ON DELETE RESTRICT;
