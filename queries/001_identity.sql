CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- COUNTRIES
-- ============================================================
CREATE TABLE IF NOT EXISTS identity_countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(450) NOT NULL,
    code3 VARCHAR(3),
    is_billing_enabled BOOLEAN DEFAULT true,
    is_shipping_enabled BOOLEAN DEFAULT true,
    is_city_enabled BOOLEAN DEFAULT true,
    is_zip_code_enabled BOOLEAN DEFAULT true,
    is_district_enabled BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE identity_countries IS 'List of countries for addresses';

-- ============================================================
-- ROLES
-- ============================================================
CREATE TABLE IF NOT EXISTS identity_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    is_system BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE identity_roles IS 'User roles for authorization';

-- ============================================================
-- VENDORS
-- ============================================================
CREATE TABLE IF NOT EXISTS identity_vendors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    email VARCHAR(255),
    phone VARCHAR(50),
    website VARCHAR(500),
    logo_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE identity_vendors IS 'Vendors/sellers in the marketplace';

CREATE INDEX IF NOT EXISTS idx_identity_vendors_slug ON identity_vendors(slug);

-- ============================================================
-- USERS
-- ============================================================
CREATE TABLE IF NOT EXISTS identity_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id UUID REFERENCES identity_vendors(id) ON DELETE RESTRICT,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50),
    password_hash VARCHAR(500) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    last_login_at TIMESTAMPTZ,
    default_shipping_address_id UUID,
    default_billing_address_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE identity_users IS 'Registered users of the system';

CREATE INDEX IF NOT EXISTS idx_identity_users_vendor_id ON identity_users(vendor_id);
CREATE INDEX IF NOT EXISTS idx_identity_users_email ON identity_users(email);
CREATE INDEX IF NOT EXISTS idx_identity_users_default_shipping_address_id ON identity_users(default_shipping_address_id);
CREATE INDEX IF NOT EXISTS idx_identity_users_default_billing_address_id ON identity_users(default_billing_address_id);

-- ============================================================
-- USER ROLES (composite PK)
-- ============================================================
CREATE TABLE IF NOT EXISTS identity_user_roles (
    user_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
    role_id UUID NOT NULL REFERENCES identity_roles(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

COMMENT ON TABLE identity_user_roles IS 'Many-to-many relationship between users and roles';

CREATE INDEX IF NOT EXISTS idx_identity_user_roles_role_id ON identity_user_roles(role_id);

-- ============================================================
-- STATE OR PROVINCES
-- ============================================================
CREATE TABLE IF NOT EXISTS identity_state_or_provinces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_id UUID NOT NULL REFERENCES identity_countries(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE identity_state_or_provinces IS 'States or provinces within a country';

CREATE INDEX IF NOT EXISTS idx_identity_state_or_provinces_country_id ON identity_state_or_provinces(country_id);

-- ============================================================
-- DISTRICTS
-- ============================================================
CREATE TABLE IF NOT EXISTS identity_districts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    state_or_province_id UUID NOT NULL REFERENCES identity_state_or_provinces(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE identity_districts IS 'Districts within a state or province';

CREATE INDEX IF NOT EXISTS idx_identity_districts_state_or_province_id ON identity_districts(state_or_province_id);

-- ============================================================
-- ADDRESSES
-- ============================================================
CREATE TABLE IF NOT EXISTS identity_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_id UUID NOT NULL REFERENCES identity_countries(id) ON DELETE RESTRICT,
    state_or_province_id UUID REFERENCES identity_state_or_provinces(id) ON DELETE RESTRICT,
    district_id UUID REFERENCES identity_districts(id) ON DELETE RESTRICT,
    city VARCHAR(255),
    zip_code VARCHAR(20),
    address_line1 VARCHAR(500) NOT NULL,
    address_line2 VARCHAR(500),
    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),
    phone VARCHAR(50),
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE identity_addresses IS 'Addresses for users and organizations';

CREATE INDEX IF NOT EXISTS idx_identity_addresses_country_id ON identity_addresses(country_id);
CREATE INDEX IF NOT EXISTS idx_identity_addresses_state_or_province_id ON identity_addresses(state_or_province_id);
CREATE INDEX IF NOT EXISTS idx_identity_addresses_district_id ON identity_addresses(district_id);

-- ============================================================
-- USER ADDRESSES
-- ============================================================
CREATE TABLE IF NOT EXISTS identity_user_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
    address_id UUID NOT NULL REFERENCES identity_addresses(id) ON DELETE RESTRICT,
    address_type VARCHAR(50) NOT NULL CHECK (address_type IN ('billing', 'shipping', 'both')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, address_id)
);

COMMENT ON TABLE identity_user_addresses IS 'Links users to their addresses with type classification';

CREATE INDEX IF NOT EXISTS idx_identity_user_addresses_user_id ON identity_user_addresses(user_id);
CREATE INDEX IF NOT EXISTS idx_identity_user_addresses_address_id ON identity_user_addresses(address_id);

-- Add FKs for default addresses on users (after addresses table exists)
ALTER TABLE identity_users ADD CONSTRAINT fk_users_default_shipping_address
    FOREIGN KEY (default_shipping_address_id) REFERENCES identity_addresses(id) ON DELETE SET NULL;

ALTER TABLE identity_users ADD CONSTRAINT fk_users_default_billing_address
    FOREIGN KEY (default_billing_address_id) REFERENCES identity_addresses(id) ON DELETE SET NULL;

-- ============================================================
-- CUSTOMER GROUPS
-- ============================================================
CREATE TABLE IF NOT EXISTS identity_customer_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE identity_customer_groups IS 'Customer groups for segmentation and pricing';

-- ============================================================
-- CUSTOMER GROUP USERS (composite PK)
-- ============================================================
CREATE TABLE IF NOT EXISTS identity_customer_group_users (
    customer_group_id UUID NOT NULL REFERENCES identity_customer_groups(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (customer_group_id, user_id)
);

COMMENT ON TABLE identity_customer_group_users IS 'Many-to-many between customer groups and users';

CREATE INDEX IF NOT EXISTS idx_identity_customer_group_users_user_id ON identity_customer_group_users(user_id);

-- ============================================================
-- MEDIA
-- ============================================================
CREATE TABLE IF NOT EXISTS identity_media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES identity_users(id) ON DELETE RESTRICT,
    file_name VARCHAR(500) NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT,
    mime_type VARCHAR(255),
    alt_text VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE identity_media IS 'Media/files uploaded by users';

CREATE INDEX IF NOT EXISTS idx_identity_media_user_id ON identity_media(user_id);
