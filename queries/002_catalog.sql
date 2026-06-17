-- ============================================================
-- CATEGORIES
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id UUID REFERENCES catalog_categories(id) ON DELETE RESTRICT,
    thumbnail_image_id UUID,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    seo_title VARCHAR(255),
    seo_description TEXT,
    is_active BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_categories IS 'Product categories with self-referencing parent hierarchy';

CREATE INDEX IF NOT EXISTS idx_catalog_categories_parent_id ON catalog_categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_catalog_categories_slug ON catalog_categories(slug);

-- ============================================================
-- BRANDS
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_brands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    logo_url TEXT,
    website VARCHAR(500),
    seo_title VARCHAR(255),
    seo_description TEXT,
    is_active BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_brands IS 'Product brands';

CREATE INDEX IF NOT EXISTS idx_catalog_brands_slug ON catalog_brands(slug);

-- ============================================================
-- MEDIA
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_name VARCHAR(500) NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT,
    mime_type VARCHAR(255),
    alt_text VARCHAR(500),
    width INT,
    height INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_media IS 'Media assets for catalog entities';

-- ============================================================
-- PRODUCTS
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand_id UUID REFERENCES catalog_brands(id) ON DELETE RESTRICT,
    tax_class_id UUID,
    thumbnail_image_id UUID REFERENCES catalog_media(id) ON DELETE RESTRICT,
    vendor_id UUID REFERENCES identity_vendors(id) ON DELETE RESTRICT,
    sku VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(500) NOT NULL,
    slug VARCHAR(500) NOT NULL UNIQUE,
    short_description TEXT,
    description TEXT,
    seo_title VARCHAR(255),
    seo_description TEXT,
    price DECIMAL(18,2) NOT NULL CHECK (price >= 0),
    compare_at_price DECIMAL(18,2) CHECK (compare_at_price >= 0),
    cost_price DECIMAL(18,2) CHECK (cost_price >= 0),
    weight DECIMAL(10,2) CHECK (weight >= 0),
    height DECIMAL(10,2) CHECK (height >= 0),
    width DECIMAL(10,2) CHECK (width >= 0),
    length DECIMAL(10,2) CHECK (length >= 0),
    is_active BOOLEAN DEFAULT true,
    is_featured BOOLEAN DEFAULT false,
    is_taxable BOOLEAN DEFAULT true,
    is_free_shipping BOOLEAN DEFAULT false,
    meta_data JSONB,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_products IS 'Main product catalog table';

CREATE INDEX IF NOT EXISTS idx_catalog_products_brand_id ON catalog_products(brand_id);
CREATE INDEX IF NOT EXISTS idx_catalog_products_tax_class_id ON catalog_products(tax_class_id);
CREATE INDEX IF NOT EXISTS idx_catalog_products_thumbnail_image_id ON catalog_products(thumbnail_image_id);
CREATE INDEX IF NOT EXISTS idx_catalog_products_vendor_id ON catalog_products(vendor_id);
CREATE INDEX IF NOT EXISTS idx_catalog_products_slug ON catalog_products(slug);
CREATE INDEX IF NOT EXISTS idx_catalog_products_sku ON catalog_products(sku);

-- ============================================================
-- PRODUCT CATEGORIES
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_product_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    category_id UUID NOT NULL REFERENCES catalog_categories(id) ON DELETE RESTRICT,
    is_featured_product BOOLEAN DEFAULT false,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_product_categories IS 'Many-to-many between products and categories';

CREATE INDEX IF NOT EXISTS idx_catalog_product_categories_product_id ON catalog_product_categories(product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_categories_category_id ON catalog_product_categories(category_id);

-- ============================================================
-- PRODUCT MEDIAS
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_product_medias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    media_id UUID NOT NULL REFERENCES catalog_media(id) ON DELETE RESTRICT,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_product_medias IS 'Many-to-many between products and media';

CREATE INDEX IF NOT EXISTS idx_catalog_product_medias_product_id ON catalog_product_medias(product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_medias_media_id ON catalog_product_medias(media_id);

-- ============================================================
-- PRODUCT ATTRIBUTE GROUPS
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_product_attribute_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_product_attribute_groups IS 'Groups for product attributes';

-- ============================================================
-- PRODUCT ATTRIBUTES
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_product_attributes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID REFERENCES catalog_product_attribute_groups(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    attribute_type VARCHAR(50) NOT NULL DEFAULT 'text',
    is_required BOOLEAN DEFAULT false,
    is_filterable BOOLEAN DEFAULT false,
    is_visible_on_product_page BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_product_attributes IS 'Product attribute definitions';

CREATE INDEX IF NOT EXISTS idx_catalog_product_attributes_group_id ON catalog_product_attributes(group_id);

-- ============================================================
-- PRODUCT ATTRIBUTE VALUES
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_product_attribute_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attribute_id UUID NOT NULL REFERENCES catalog_product_attributes(id) ON DELETE RESTRICT,
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    value_text TEXT,
    value_boolean BOOLEAN,
    value_number DECIMAL(18,4),
    value_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_product_attribute_values IS 'Values of attributes for specific products';

CREATE INDEX IF NOT EXISTS idx_catalog_product_attribute_values_attribute_id ON catalog_product_attribute_values(attribute_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_attribute_values_product_id ON catalog_product_attribute_values(product_id);

-- ============================================================
-- PRODUCT OPTIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_product_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_product_options IS 'Product options (e.g., size, color)';

-- ============================================================
-- PRODUCT OPTION VALUES
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_product_option_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    option_id UUID NOT NULL REFERENCES catalog_product_options(id) ON DELETE RESTRICT,
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_product_option_values IS 'Values for product options linked to products';

CREATE INDEX IF NOT EXISTS idx_catalog_product_option_values_option_id ON catalog_product_option_values(option_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_option_values_product_id ON catalog_product_option_values(product_id);

-- ============================================================
-- PRODUCT OPTION COMBINATIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_product_option_combinations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    option_id UUID NOT NULL REFERENCES catalog_product_options(id) ON DELETE RESTRICT,
    sku VARCHAR(255),
    price_adjustment DECIMAL(18,2) DEFAULT 0 CHECK (price_adjustment >= 0),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_product_option_combinations IS 'Links product options to products with pricing';

CREATE INDEX IF NOT EXISTS idx_catalog_product_option_combinations_product_id ON catalog_product_option_combinations(product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_option_combinations_option_id ON catalog_product_option_combinations(option_id);

-- ============================================================
-- PRODUCT LINKS
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_product_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    linked_product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    link_type VARCHAR(50) NOT NULL CHECK (link_type IN ('related', 'upsell', 'cross_sell', 'bundle')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_product_links IS 'Links between related products';

CREATE INDEX IF NOT EXISTS idx_catalog_product_links_product_id ON catalog_product_links(product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_links_linked_product_id ON catalog_product_links(linked_product_id);

-- ============================================================
-- PRODUCT TEMPLATES
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_product_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_product_templates IS 'Product templates for consistent attribute sets';

-- ============================================================
-- PRODUCT TEMPLATE PRODUCT ATTRIBUTES (composite PK)
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_product_template_product_attributes (
    product_template_id UUID NOT NULL REFERENCES catalog_product_templates(id) ON DELETE RESTRICT,
    product_attribute_id UUID NOT NULL REFERENCES catalog_product_attributes(id) ON DELETE RESTRICT,
    is_required BOOLEAN DEFAULT false,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (product_template_id, product_attribute_id)
);

COMMENT ON TABLE catalog_product_template_product_attributes IS 'Many-to-many between templates and attributes';

CREATE INDEX IF NOT EXISTS idx_catalog_product_template_product_attributes_attr_id ON catalog_product_template_product_attributes(product_attribute_id);

-- ============================================================
-- PRODUCT PRICE HISTORIES
-- ============================================================
CREATE TABLE IF NOT EXISTS catalog_product_price_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog_products(id) ON DELETE RESTRICT,
    old_price DECIMAL(18,2) CHECK (old_price >= 0),
    new_price DECIMAL(18,2) NOT NULL CHECK (new_price >= 0),
    created_by_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE catalog_product_price_histories IS 'Audit log of product price changes';

CREATE INDEX IF NOT EXISTS idx_catalog_product_price_histories_product_id ON catalog_product_price_histories(product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_price_histories_created_by_id ON catalog_product_price_histories(created_by_id);

-- Add FK from catalog_categories.thumbnail_image_id to catalog_media
ALTER TABLE catalog_categories ADD CONSTRAINT fk_categories_thumbnail_image
    FOREIGN KEY (thumbnail_image_id) REFERENCES catalog_media(id) ON DELETE SET NULL;
