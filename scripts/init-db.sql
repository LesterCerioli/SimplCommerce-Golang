-- ============================================================================
-- SimpleCommerce-Go Database Schema
-- PostgreSQL Migration Script
-- ============================================================================
-- This script creates all tables for the SimplCommerce-Go commerce platform.
-- Naming convention: {module}_{entity} (lowercase plural)
-- All tables use BIGSERIAL primary keys, TIMESTAMPTZ for dates,
-- DECIMAL(18,2) for monetary values, and VARCHAR for strings.
-- ============================================================================

-- Enable UUID support (used for user GUIDs and external references)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- IDENTITY MODULE
-- ============================================================================

-- Users of the system (customers, admins, etc.)
CREATE TABLE IF NOT EXISTS identity_users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_guid VARCHAR(36) NOT NULL,
    full_name VARCHAR(450) NOT NULL,
    email VARCHAR(256) NOT NULL,
    password_hash VARCHAR(500) NOT NULL,
    phone_number VARCHAR(50),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    culture VARCHAR(10),
    refresh_token_hash VARCHAR(500),
    vendor_id BIGINT,
    default_shipping_address_id BIGINT,
    default_billing_address_id BIGINT
);

COMMENT ON TABLE identity_users IS 'System users including customers, administrators, and vendors';

-- Roles for role-based access control
CREATE TABLE IF NOT EXISTS identity_roles (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(256) NOT NULL,
    normalized_name VARCHAR(256) NOT NULL,
    concurrency_stamp VARCHAR(256)
);

COMMENT ON TABLE identity_roles IS 'Application roles for authorization';

-- Join table: user-role assignments (many-to-many)
CREATE TABLE IF NOT EXISTS identity_user_roles (
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, role_id)
);

COMMENT ON TABLE identity_user_roles IS 'Many-to-many relationship between users and roles';

-- Addresses
CREATE TABLE IF NOT EXISTS identity_addresses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    contact_name VARCHAR(450) NOT NULL,
    phone VARCHAR(50),
    address_line1 VARCHAR(450) NOT NULL,
    address_line2 VARCHAR(450),
    city VARCHAR(200) NOT NULL,
    zip_code VARCHAR(20),
    district_id BIGINT,
    state_or_province_id BIGINT NOT NULL,
    country_id VARCHAR(10)
);

COMMENT ON TABLE identity_addresses IS 'Physical addresses for users and organizations';

-- User-address assignments (a user can have multiple addresses)
CREATE TABLE IF NOT EXISTS identity_user_addresses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,
    address_id BIGINT NOT NULL,
    address_type INT NOT NULL DEFAULT 0
);

COMMENT ON TABLE identity_user_addresses IS 'Maps users to their addresses with type (shipping/billing)';

-- Countries
CREATE TABLE IF NOT EXISTS identity_countries (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(450) NOT NULL,
    code3 VARCHAR(3),
    is_billing_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    is_shipping_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    is_city_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    is_zip_code_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    is_district_enabled BOOLEAN NOT NULL DEFAULT FALSE
);

COMMENT ON TABLE identity_countries IS 'Countries for address and tax configuration';

-- States or provinces within a country
CREATE TABLE IF NOT EXISTS identity_state_or_provinces (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    country_id VARCHAR(10) NOT NULL,
    code VARCHAR(10),
    name VARCHAR(450) NOT NULL,
    type VARCHAR(50)
);

COMMENT ON TABLE identity_state_or_provinces IS 'States or provinces within a country';

-- Districts within a state/province
CREATE TABLE IF NOT EXISTS identity_districts (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    state_or_province_id BIGINT NOT NULL,
    name VARCHAR(450) NOT NULL,
    type VARCHAR(50),
    location TEXT
);

COMMENT ON TABLE identity_districts IS 'Districts within a state or province';

-- Vendors (suppliers/sellers)
CREATE TABLE IF NOT EXISTS identity_vendors (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    slug VARCHAR(450) NOT NULL,
    description TEXT,
    email VARCHAR(256),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

COMMENT ON TABLE identity_vendors IS 'Product vendors/sellers';

-- Customer groups for segmentation and pricing rules
CREATE TABLE IF NOT EXISTS identity_customer_groups (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

COMMENT ON TABLE identity_customer_groups IS 'Customer groups for targeting promotions and pricing';

-- Join table: customer group-user assignments
CREATE TABLE IF NOT EXISTS identity_customer_group_users (
    user_id BIGINT NOT NULL,
    customer_group_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, customer_group_id)
);

COMMENT ON TABLE identity_customer_group_users IS 'Many-to-many relationship between users and customer groups';

-- Media (images, videos, files) - identity module
CREATE TABLE IF NOT EXISTS identity_media (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    caption VARCHAR(450),
    file_size BIGINT NOT NULL DEFAULT 0,
    file_name VARCHAR(450) NOT NULL,
    media_type INT NOT NULL DEFAULT 0
);

COMMENT ON TABLE identity_media IS 'Media files (images, videos, documents) attached to identity entities';

-- Widget zones (layout regions for widgets)
CREATE TABLE IF NOT EXISTS identity_widget_zones (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    description TEXT
);

COMMENT ON TABLE identity_widget_zones IS 'Widget zones/layout regions for page rendering';

-- Widgets (available widget types)
CREATE TABLE IF NOT EXISTS identity_widgets (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    view_component_name VARCHAR(450),
    is_system BOOLEAN NOT NULL DEFAULT FALSE
);

COMMENT ON TABLE identity_widgets IS 'Available widget types in the system';

-- Widget instances (specific widget placements in zones)
CREATE TABLE IF NOT EXISTS identity_widget_instances (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    widget_id BIGINT NOT NULL,
    widget_zone_id BIGINT NOT NULL,
    display_order INT NOT NULL DEFAULT 0,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    data TEXT
);

COMMENT ON TABLE identity_widget_instances IS 'Instances of widgets placed in widget zones';

-- ============================================================================
-- CATALOG MODULE
-- ============================================================================

-- Products
CREATE TABLE IF NOT EXISTS catalog_products (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    slug VARCHAR(450) NOT NULL,
    short_description VARCHAR(450),
    description TEXT,
    specification TEXT,
    price DECIMAL(18,2) NOT NULL DEFAULT 0,
    old_price DECIMAL(18,2),
    special_price DECIMAL(18,2),
    special_price_start TIMESTAMPTZ,
    special_price_end TIMESTAMPTZ,
    has_options BOOLEAN NOT NULL DEFAULT FALSE,
    is_visible_individually BOOLEAN NOT NULL DEFAULT TRUE,
    is_featured BOOLEAN NOT NULL DEFAULT FALSE,
    is_call_for_pricing BOOLEAN NOT NULL DEFAULT FALSE,
    is_allow_to_order BOOLEAN NOT NULL DEFAULT TRUE,
    stock_tracking_is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    stock_quantity INT NOT NULL DEFAULT 0,
    sku VARCHAR(200),
    gtin VARCHAR(50),
    normalized_name VARCHAR(450),
    display_order INT NOT NULL DEFAULT 0,
    reviews_count INT NOT NULL DEFAULT 0,
    rating_average DECIMAL(3,2),
    vendor_id BIGINT,
    brand_id BIGINT,
    tax_class_id BIGINT,
    thumbnail_image_id BIGINT,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    published_on TIMESTAMPTZ
);

COMMENT ON TABLE catalog_products IS 'Core product catalog items';

-- Categories
CREATE TABLE IF NOT EXISTS catalog_categories (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    slug VARCHAR(450) NOT NULL,
    description TEXT,
    display_order INT NOT NULL DEFAULT 0,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    include_in_menu BOOLEAN NOT NULL DEFAULT TRUE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    parent_id BIGINT,
    thumbnail_image_id BIGINT
);

COMMENT ON TABLE catalog_categories IS 'Product categories with hierarchical parent-child relationships';

-- Brands
CREATE TABLE IF NOT EXISTS catalog_brands (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    slug VARCHAR(450) NOT NULL,
    description TEXT,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

COMMENT ON TABLE catalog_brands IS 'Product brands';

-- Join table: product-category assignments with metadata
CREATE TABLE IF NOT EXISTS catalog_product_categories (
    product_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    is_featured_product BOOLEAN NOT NULL DEFAULT FALSE,
    display_order INT NOT NULL DEFAULT 0,
    PRIMARY KEY (product_id, category_id)
);

COMMENT ON TABLE catalog_product_categories IS 'Many-to-many relationship between products and categories';

-- Product attribute groups
CREATE TABLE IF NOT EXISTS catalog_product_attribute_groups (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL
);

COMMENT ON TABLE catalog_product_attribute_groups IS 'Groups for organizing product attributes';

-- Product attributes
CREATE TABLE IF NOT EXISTS catalog_product_attributes (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    group_id BIGINT
);

COMMENT ON TABLE catalog_product_attributes IS 'Product attribute definitions (e.g., Color, Size)';

-- Product attribute values (actual values assigned to products)
CREATE TABLE IF NOT EXISTS catalog_product_attribute_values (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    attribute_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    value TEXT
);

COMMENT ON TABLE catalog_product_attribute_values IS 'Values of attributes assigned to specific products';

-- Product options (e.g., Size, Color choices)
CREATE TABLE IF NOT EXISTS catalog_product_options (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL
);

COMMENT ON TABLE catalog_product_options IS 'Product option definitions for configurable products';

-- Product option values
CREATE TABLE IF NOT EXISTS catalog_product_option_values (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    option_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    value VARCHAR(450),
    display_type VARCHAR(50),
    sort_index INT NOT NULL DEFAULT 0
);

COMMENT ON TABLE catalog_product_option_values IS 'Values of options assigned to specific products';

-- Product option combinations (for variant tracking)
CREATE TABLE IF NOT EXISTS catalog_product_option_combinations (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    product_id BIGINT NOT NULL,
    option_id BIGINT NOT NULL,
    value VARCHAR(450),
    sort_index INT NOT NULL DEFAULT 0
);

COMMENT ON TABLE catalog_product_option_combinations IS 'Combinations of product options for variant management';

-- Product links (related, cross-sell, up-sell, super)
CREATE TABLE IF NOT EXISTS catalog_product_links (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    product_id BIGINT NOT NULL,
    linked_product_id BIGINT NOT NULL,
    link_type INT NOT NULL DEFAULT 2
);

COMMENT ON TABLE catalog_product_links IS 'Linked products (related, cross-sell, up-sell, super)';

-- Product media assignments
CREATE TABLE IF NOT EXISTS catalog_product_medias (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    product_id BIGINT NOT NULL,
    media_id BIGINT NOT NULL,
    display_order INT NOT NULL DEFAULT 0
);

COMMENT ON TABLE catalog_product_medias IS 'Media attachments for products';

-- Product templates
CREATE TABLE IF NOT EXISTS catalog_product_templates (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL
);

COMMENT ON TABLE catalog_product_templates IS 'Product templates defining attribute sets';

-- Join table: template-attribute assignments
CREATE TABLE IF NOT EXISTS catalog_product_template_product_attributes (
    product_template_id BIGINT NOT NULL,
    product_attribute_id BIGINT NOT NULL,
    PRIMARY KEY (product_template_id, product_attribute_id)
);

COMMENT ON TABLE catalog_product_template_product_attributes IS 'Many-to-many relationship between templates and attributes';

-- Product price change history
CREATE TABLE IF NOT EXISTS catalog_product_price_histories (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    product_id BIGINT NOT NULL,
    created_by_id BIGINT NOT NULL DEFAULT 0,
    price DECIMAL(18,2),
    old_price DECIMAL(18,2),
    special_price DECIMAL(18,2),
    special_price_start TIMESTAMPTZ,
    special_price_end TIMESTAMPTZ
);

COMMENT ON TABLE catalog_product_price_histories IS 'Audit log of product price changes';

-- Catalog media
CREATE TABLE IF NOT EXISTS catalog_media (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    caption VARCHAR(450),
    file_size BIGINT NOT NULL DEFAULT 0,
    file_name VARCHAR(450) NOT NULL,
    media_type INT NOT NULL DEFAULT 0
);

COMMENT ON TABLE catalog_media IS 'Media files attached to catalog entities';

-- ============================================================================
-- CART MODULE
-- ============================================================================

-- Shopping cart items
CREATE TABLE IF NOT EXISTS cart_cart_items (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    product_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    vendor_id BIGINT,
    CONSTRAINT cart_items_quantity_check CHECK (quantity >= 1)
);

COMMENT ON TABLE cart_cart_items IS 'Shopping cart items for each customer';

-- ============================================================================
-- ORDERS MODULE
-- ============================================================================

-- Orders
CREATE TABLE IF NOT EXISTS orders_orders (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    customer_id BIGINT NOT NULL,
    vendor_id BIGINT,
    created_by_id BIGINT NOT NULL DEFAULT 0,
    updated_by_id BIGINT,
    coupon_code VARCHAR(100),
    coupon_rule_name VARCHAR(450),
    discount_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    sub_total DECIMAL(18,2) NOT NULL,
    sub_total_with_discount DECIMAL(18,2) NOT NULL,
    shipping_address_id BIGINT NOT NULL,
    billing_address_id BIGINT NOT NULL,
    order_status VARCHAR(50) NOT NULL DEFAULT 'New',
    order_note VARCHAR(1000),
    parent_id BIGINT,
    is_master_order BOOLEAN NOT NULL DEFAULT FALSE,
    shipping_method VARCHAR(450),
    shipping_fee_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    order_total DECIMAL(18,2) NOT NULL,
    payment_method VARCHAR(450),
    payment_fee_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    CONSTRAINT orders_sub_total_check CHECK (sub_total >= 0),
    CONSTRAINT orders_total_check CHECK (order_total >= 0)
);

COMMENT ON TABLE orders_orders IS 'Customer orders';

-- Order line items
CREATE TABLE IF NOT EXISTS orders_order_items (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    product_name VARCHAR(450),
    product_sku VARCHAR(200),
    product_price DECIMAL(18,2) NOT NULL,
    quantity INT NOT NULL,
    discount_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    tax_percent DECIMAL(18,2) NOT NULL DEFAULT 0,
    CONSTRAINT order_items_quantity_check CHECK (quantity >= 1)
);

COMMENT ON TABLE orders_order_items IS 'Individual line items within an order';

-- Order addresses (snapshot of address at order time)
CREATE TABLE IF NOT EXISTS orders_order_addresses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    contact_name VARCHAR(450) NOT NULL,
    phone VARCHAR(50),
    address_line1 VARCHAR(450) NOT NULL,
    address_line2 VARCHAR(450),
    city VARCHAR(200) NOT NULL,
    zip_code VARCHAR(20),
    district_id BIGINT,
    state_or_province_id BIGINT NOT NULL,
    country_id VARCHAR(10) NOT NULL
);

COMMENT ON TABLE orders_order_addresses IS 'Snapshot of shipping/billing address at time of order';

-- Order status history
CREATE TABLE IF NOT EXISTS orders_order_histories (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    order_id BIGINT NOT NULL,
    old_status VARCHAR(50),
    new_status VARCHAR(50) NOT NULL,
    note VARCHAR(1000),
    created_by_id BIGINT NOT NULL DEFAULT 0
);

COMMENT ON TABLE orders_order_histories IS 'Audit trail of order status changes';

-- ============================================================================
-- PAYMENTS MODULE
-- ============================================================================

-- Payments
CREATE TABLE IF NOT EXISTS payments_payments (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    order_id BIGINT NOT NULL,
    amount DECIMAL(18,2) NOT NULL,
    payment_fee DECIMAL(18,2) NOT NULL DEFAULT 0,
    payment_method VARCHAR(450) NOT NULL,
    gateway_transaction_id VARCHAR(500),
    status VARCHAR(50) NOT NULL DEFAULT 'Succeeded',
    failure_message TEXT
);

COMMENT ON TABLE payments_payments IS 'Payment transactions for orders';

-- Payment providers (integration configurations)
CREATE TABLE IF NOT EXISTS payments_payment_providers (
    id VARCHAR(200) PRIMARY KEY,
    name VARCHAR(450) NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    configure_url VARCHAR(450),
    landing_view_component_name VARCHAR(450),
    additional_settings TEXT
);

COMMENT ON TABLE payments_payment_providers IS 'Configured payment gateway providers';

-- ============================================================================
-- SHIPPING MODULE
-- ============================================================================

-- Shipments
CREATE TABLE IF NOT EXISTS shipping_shipments (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    order_id BIGINT NOT NULL,
    tracking_number VARCHAR(450),
    warehouse_id BIGINT NOT NULL,
    vendor_id BIGINT,
    created_by_id BIGINT NOT NULL DEFAULT 0
);

COMMENT ON TABLE shipping_shipments IS 'Order shipments';

-- Shipment items
CREATE TABLE IF NOT EXISTS shipping_shipment_items (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    shipment_id BIGINT NOT NULL,
    order_item_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    CONSTRAINT shipment_items_quantity_check CHECK (quantity >= 1)
);

COMMENT ON TABLE shipping_shipment_items IS 'Individual items within a shipment';

-- Shipping providers
CREATE TABLE IF NOT EXISTS shipping_shipping_providers (
    id VARCHAR(200) PRIMARY KEY,
    name VARCHAR(450) NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    configure_url VARCHAR(450),
    to_all_shipping_enabled_countries BOOLEAN NOT NULL DEFAULT TRUE,
    only_country_ids_string VARCHAR(1000),
    to_all_shipping_enabled_states_or_provinces BOOLEAN NOT NULL DEFAULT TRUE,
    only_state_or_province_ids_string VARCHAR(1000),
    additional_settings TEXT,
    shipping_price_service_type_name VARCHAR(450)
);

COMMENT ON TABLE shipping_shipping_providers IS 'Configured shipping providers';

-- Table rate shipping prices and destinations
CREATE TABLE IF NOT EXISTS shipping_table_rate_price_and_destinations (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    country_id VARCHAR(10),
    state_or_province_id BIGINT,
    district_id BIGINT,
    zip_code VARCHAR(20),
    note TEXT,
    min_order_subtotal DECIMAL(18,2) NOT NULL DEFAULT 0,
    shipping_price DECIMAL(18,2) NOT NULL
);

COMMENT ON TABLE shipping_table_rate_price_and_destinations IS 'Shipping rates based on destination and order subtotal';

-- ============================================================================
-- INVENTORY MODULE
-- ============================================================================

-- Stock levels per product per warehouse
CREATE TABLE IF NOT EXISTS inventory_stocks (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    product_id BIGINT NOT NULL,
    warehouse_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    reserved_quantity INT NOT NULL DEFAULT 0,
    CONSTRAINT stocks_quantity_check CHECK (quantity >= 0),
    CONSTRAINT stocks_reserved_check CHECK (reserved_quantity >= 0),
    CONSTRAINT stocks_reserved_not_exceed CHECK (reserved_quantity <= quantity)
);

COMMENT ON TABLE inventory_stocks IS 'Product inventory stock levels per warehouse';

-- Stock adjustment history
CREATE TABLE IF NOT EXISTS inventory_stock_histories (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    product_id BIGINT NOT NULL,
    warehouse_id BIGINT NOT NULL,
    created_by_id BIGINT NOT NULL DEFAULT 0,
    adjusted_quantity INT NOT NULL DEFAULT 0,
    note VARCHAR(1000)
);

COMMENT ON TABLE inventory_stock_histories IS 'Audit trail of stock adjustments';

-- Warehouses
CREATE TABLE IF NOT EXISTS inventory_warehouses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    address_id BIGINT NOT NULL,
    vendor_id BIGINT
);

COMMENT ON TABLE inventory_warehouses IS 'Physical warehouse locations';

-- Back-in-stock subscriptions
CREATE TABLE IF NOT EXISTS inventory_product_back_in_stock_subscriptions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    product_id BIGINT NOT NULL,
    customer_email VARCHAR(256) NOT NULL
);

COMMENT ON TABLE inventory_product_back_in_stock_subscriptions IS 'Customer email subscriptions for out-of-stock product notifications';

-- ============================================================================
-- PRICING MODULE
-- ============================================================================

-- Cart price rules (discount rules for shopping cart)
CREATE TABLE IF NOT EXISTS pricing_cart_rules (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    start_on TIMESTAMPTZ,
    end_on TIMESTAMPTZ,
    is_coupon_required BOOLEAN NOT NULL DEFAULT FALSE,
    rule_to_apply VARCHAR(50) NOT NULL DEFAULT 'by_fixed',
    discount_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    max_discount_amount DECIMAL(18,2),
    discount_step INT,
    usage_limit_per_coupon INT,
    usage_limit_per_customer INT
);

COMMENT ON TABLE pricing_cart_rules IS 'Cart-level discount and promotion rules';

-- Coupons
CREATE TABLE IF NOT EXISTS pricing_coupons (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    cart_rule_id BIGINT NOT NULL,
    code VARCHAR(100) NOT NULL
);

COMMENT ON TABLE pricing_coupons IS 'Coupon codes associated with cart rules';

-- Join table: cart rule-category assignments
CREATE TABLE IF NOT EXISTS pricing_cart_rule_categories (
    cart_rule_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    PRIMARY KEY (cart_rule_id, category_id)
);

COMMENT ON TABLE pricing_cart_rule_categories IS 'Categories that a cart rule applies to';

-- Join table: cart rule-product assignments
CREATE TABLE IF NOT EXISTS pricing_cart_rule_products (
    cart_rule_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    PRIMARY KEY (cart_rule_id, product_id)
);

COMMENT ON TABLE pricing_cart_rule_products IS 'Products that a cart rule applies to';

-- Join table: cart rule-customer group assignments
CREATE TABLE IF NOT EXISTS pricing_cart_rule_customer_groups (
    cart_rule_id BIGINT NOT NULL,
    customer_group_id BIGINT NOT NULL,
    PRIMARY KEY (cart_rule_id, customer_group_id)
);

COMMENT ON TABLE pricing_cart_rule_customer_groups IS 'Customer groups that a cart rule applies to';

-- Cart rule usage tracking
CREATE TABLE IF NOT EXISTS pricing_cart_rule_usages (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    cart_rule_id BIGINT NOT NULL,
    coupon_id BIGINT,
    user_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL
);

COMMENT ON TABLE pricing_cart_rule_usages IS 'Tracks usage of cart rules and coupons per order';

-- Catalog price rules (automatic discounts on product listing)
CREATE TABLE IF NOT EXISTS pricing_catalog_rules (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    start_on TIMESTAMPTZ,
    end_on TIMESTAMPTZ,
    rule_to_apply VARCHAR(50) NOT NULL DEFAULT 'by_fixed',
    discount_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    max_discount_amount DECIMAL(18,2)
);

COMMENT ON TABLE pricing_catalog_rules IS 'Catalog-level automatic discount rules';

-- Join table: catalog rule-customer group assignments
CREATE TABLE IF NOT EXISTS pricing_catalog_rule_customer_groups (
    catalog_rule_id BIGINT NOT NULL,
    customer_group_id BIGINT NOT NULL,
    PRIMARY KEY (catalog_rule_id, customer_group_id)
);

COMMENT ON TABLE pricing_catalog_rule_customer_groups IS 'Customer groups that a catalog rule applies to';

-- ============================================================================
-- TAX MODULE
-- ============================================================================

-- Tax classes
CREATE TABLE IF NOT EXISTS tax_tax_classes (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL
);

COMMENT ON TABLE tax_tax_classes IS 'Tax class definitions (e.g., Standard, Reduced, Zero)';

-- Tax rates
CREATE TABLE IF NOT EXISTS tax_tax_rates (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    tax_class_id BIGINT NOT NULL,
    country_id VARCHAR(10) NOT NULL,
    state_or_province_id BIGINT,
    rate DECIMAL(18,2) NOT NULL,
    zip_code VARCHAR(20)
);

COMMENT ON TABLE tax_tax_rates IS 'Tax rates by tax class, country, and state';

-- ============================================================================
-- REVIEWS MODULE
-- ============================================================================

-- Product reviews
CREATE TABLE IF NOT EXISTS reviews_reviews (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,
    title VARCHAR(450),
    comment TEXT,
    rating INT NOT NULL,
    reviewer_name VARCHAR(450),
    status VARCHAR(50) NOT NULL DEFAULT 'Pending',
    entity_type_id VARCHAR(450) NOT NULL,
    entity_id BIGINT NOT NULL,
    CONSTRAINT reviews_rating_check CHECK (rating >= 1 AND rating <= 5)
);

COMMENT ON TABLE reviews_reviews IS 'Product reviews and ratings';

-- Review replies
CREATE TABLE IF NOT EXISTS reviews_replies (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    review_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    comment TEXT,
    replier_name VARCHAR(450),
    status VARCHAR(50) NOT NULL DEFAULT 'Pending'
);

COMMENT ON TABLE reviews_replies IS 'Replies to product reviews';

-- ============================================================================
-- CMS MODULE
-- ============================================================================

-- CMS pages
CREATE TABLE IF NOT EXISTS cms_pages (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    slug VARCHAR(450) NOT NULL,
    body TEXT,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    published_on TIMESTAMPTZ
);

COMMENT ON TABLE cms_pages IS 'Content management system pages';

-- Menus
CREATE TABLE IF NOT EXISTS cms_menus (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    is_system BOOLEAN NOT NULL DEFAULT FALSE
);

COMMENT ON TABLE cms_menus IS 'Navigation menus';

-- Menu items
CREATE TABLE IF NOT EXISTS cms_menu_items (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    parent_id BIGINT,
    menu_id BIGINT NOT NULL,
    entity_id BIGINT,
    custom_link VARCHAR(450),
    name VARCHAR(450) NOT NULL,
    display_order INT NOT NULL DEFAULT 0
);

COMMENT ON TABLE cms_menu_items IS 'Individual items within navigation menus';

-- ============================================================================
-- ACTIVITY LOG MODULE
-- ============================================================================

-- Activity log entries
CREATE TABLE IF NOT EXISTS activity_log_activities (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    activity_type_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    entity_id BIGINT,
    entity_type_id VARCHAR(450)
);

COMMENT ON TABLE activity_log_activities IS 'Audit log of user activities';

-- Activity types
CREATE TABLE IF NOT EXISTS activity_log_activity_types (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL
);

COMMENT ON TABLE activity_log_activity_types IS 'Types/categories of activities that can be logged';

-- ============================================================================
-- NOTIFICATIONS MODULE
-- ============================================================================

-- User notifications
CREATE TABLE IF NOT EXISTS notifications_notifications (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,
    title VARCHAR(450) NOT NULL,
    body TEXT,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    entity_id BIGINT,
    entity_type VARCHAR(450)
);

COMMENT ON TABLE notifications_notifications IS 'User notifications';

-- ============================================================================
-- SEARCH MODULE
-- ============================================================================

-- Search queries
CREATE TABLE IF NOT EXISTS search_queries (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    query_text VARCHAR(500) NOT NULL,
    results_count INT NOT NULL DEFAULT 0
);

COMMENT ON TABLE search_queries IS 'Logged search queries';

-- ============================================================================
-- COMMENTS MODULE
-- ============================================================================

-- Comments on entities (products, pages, etc.)
CREATE TABLE IF NOT EXISTS comments_comments (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,
    parent_id BIGINT,
    entity_type_id VARCHAR(450) NOT NULL,
    entity_id BIGINT NOT NULL,
    comment_text TEXT NOT NULL,
    is_approved BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(50) NOT NULL DEFAULT 'Pending'
);

COMMENT ON TABLE comments_comments IS 'User comments on entities with threading support';

-- ============================================================================
-- NEWS MODULE
-- ============================================================================

-- News/newsletter items
CREATE TABLE IF NOT EXISTS news_news_items (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    title VARCHAR(450) NOT NULL,
    slug VARCHAR(450) NOT NULL,
    short_content VARCHAR(450),
    full_content TEXT,
    thumbnail_image_id BIGINT,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    published_on TIMESTAMPTZ
);

COMMENT ON TABLE news_news_items IS 'News articles and announcements';

-- News categories
CREATE TABLE IF NOT EXISTS news_news_categories (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(450) NOT NULL,
    slug VARCHAR(450) NOT NULL,
    is_published BOOLEAN NOT NULL DEFAULT FALSE
);

COMMENT ON TABLE news_news_categories IS 'News article categories';

-- Join table: news item-category assignments
CREATE TABLE IF NOT EXISTS news_news_item_categories (
    news_item_id BIGINT NOT NULL,
    news_category_id BIGINT NOT NULL,
    PRIMARY KEY (news_item_id, news_category_id)
);

COMMENT ON TABLE news_news_item_categories IS 'Many-to-many relationship between news items and categories';

-- ============================================================================
-- PRODUCT COMPARISON MODULE
-- ============================================================================

-- Saved product comparisons
CREATE TABLE IF NOT EXISTS product_comparison_comparing_products (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL
);

COMMENT ON TABLE product_comparison_comparing_products IS 'Products saved for comparison by users';

-- ============================================================================
-- WISHLIST MODULE
-- ============================================================================

-- Wish lists
CREATE TABLE IF NOT EXISTS wish_list_wish_lists (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,
    name VARCHAR(450),
    is_shared BOOLEAN NOT NULL DEFAULT FALSE,
    sharing_code VARCHAR(36)
);

COMMENT ON TABLE wish_list_wish_lists IS 'User wish lists';

-- Wish list items
CREATE TABLE IF NOT EXISTS wish_list_wish_list_items (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    wish_list_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    CONSTRAINT wish_list_items_quantity_check CHECK (quantity >= 1)
);

COMMENT ON TABLE wish_list_wish_list_items IS 'Products within a wish list';

-- ============================================================================
-- PRODUCT RECENTLY VIEWED MODULE
-- ============================================================================

-- Recently viewed products by users
CREATE TABLE IF NOT EXISTS product_recently_viewed_recently_viewed_products (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    latest_viewed_on TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE product_recently_viewed_recently_viewed_products IS 'Products recently viewed by users';

-- ============================================================================
-- FOREIGN KEY CONSTRAINTS
-- ============================================================================

-- Identity module FKs
ALTER TABLE identity_users
    ADD CONSTRAINT fk_identity_users_vendor FOREIGN KEY (vendor_id) REFERENCES identity_vendors(id) ON DELETE RESTRICT;

ALTER TABLE identity_user_roles
    ADD CONSTRAINT fk_identity_user_roles_user FOREIGN KEY (user_id) REFERENCES identity_users(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_identity_user_roles_role FOREIGN KEY (role_id) REFERENCES identity_roles(id) ON DELETE RESTRICT;

ALTER TABLE identity_addresses
    ADD CONSTRAINT fk_identity_addresses_district FOREIGN KEY (district_id) REFERENCES identity_districts(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_identity_addresses_state_or_province FOREIGN KEY (state_or_province_id) REFERENCES identity_state_or_provinces(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_identity_addresses_country FOREIGN KEY (country_id) REFERENCES identity_countries(id) ON DELETE RESTRICT;

ALTER TABLE identity_user_addresses
    ADD CONSTRAINT fk_identity_user_addresses_user FOREIGN KEY (user_id) REFERENCES identity_users(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_identity_user_addresses_address FOREIGN KEY (address_id) REFERENCES identity_addresses(id) ON DELETE RESTRICT;

ALTER TABLE identity_state_or_provinces
    ADD CONSTRAINT fk_identity_state_or_provinces_country FOREIGN KEY (country_id) REFERENCES identity_countries(id) ON DELETE RESTRICT;

ALTER TABLE identity_districts
    ADD CONSTRAINT fk_identity_districts_state_or_province FOREIGN KEY (state_or_province_id) REFERENCES identity_state_or_provinces(id) ON DELETE RESTRICT;

ALTER TABLE identity_customer_group_users
    ADD CONSTRAINT fk_identity_customer_group_users_user FOREIGN KEY (user_id) REFERENCES identity_users(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_identity_customer_group_users_group FOREIGN KEY (customer_group_id) REFERENCES identity_customer_groups(id) ON DELETE RESTRICT;

ALTER TABLE identity_widget_instances
    ADD CONSTRAINT fk_identity_widget_instances_widget FOREIGN KEY (widget_id) REFERENCES identity_widgets(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_identity_widget_instances_zone FOREIGN KEY (widget_zone_id) REFERENCES identity_widget_zones(id) ON DELETE RESTRICT;

-- Catalog module FKs
ALTER TABLE catalog_products
    ADD CONSTRAINT fk_catalog_products_vendor FOREIGN KEY (vendor_id) REFERENCES identity_vendors(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_catalog_products_brand FOREIGN KEY (brand_id) REFERENCES catalog_brands(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_catalog_products_tax_class FOREIGN KEY (tax_class_id) REFERENCES tax_tax_classes(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_catalog_products_thumbnail FOREIGN KEY (thumbnail_image_id) REFERENCES catalog_media(id) ON DELETE RESTRICT;

ALTER TABLE catalog_categories
    ADD CONSTRAINT fk_catalog_categories_parent FOREIGN KEY (parent_id) REFERENCES catalog_categories(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_catalog_categories_thumbnail FOREIGN KEY (thumbnail_image_id) REFERENCES catalog_media(id) ON DELETE RESTRICT;

ALTER TABLE catalog_product_categories
    ADD CONSTRAINT fk_catalog_product_categories_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_catalog_product_categories_category FOREIGN KEY (category_id) REFERENCES catalog_categories(id) ON DELETE RESTRICT;

ALTER TABLE catalog_product_attributes
    ADD CONSTRAINT fk_catalog_product_attributes_group FOREIGN KEY (group_id) REFERENCES catalog_product_attribute_groups(id) ON DELETE RESTRICT;

ALTER TABLE catalog_product_attribute_values
    ADD CONSTRAINT fk_catalog_product_attribute_values_attribute FOREIGN KEY (attribute_id) REFERENCES catalog_product_attributes(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_catalog_product_attribute_values_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT;

ALTER TABLE catalog_product_option_values
    ADD CONSTRAINT fk_catalog_product_option_values_option FOREIGN KEY (option_id) REFERENCES catalog_product_options(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_catalog_product_option_values_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT;

ALTER TABLE catalog_product_option_combinations
    ADD CONSTRAINT fk_catalog_product_option_combinations_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_catalog_product_option_combinations_option FOREIGN KEY (option_id) REFERENCES catalog_product_options(id) ON DELETE RESTRICT;

ALTER TABLE catalog_product_links
    ADD CONSTRAINT fk_catalog_product_links_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_catalog_product_links_linked_product FOREIGN KEY (linked_product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT;

ALTER TABLE catalog_product_medias
    ADD CONSTRAINT fk_catalog_product_medias_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_catalog_product_medias_media FOREIGN KEY (media_id) REFERENCES catalog_media(id) ON DELETE RESTRICT;

ALTER TABLE catalog_product_template_product_attributes
    ADD CONSTRAINT fk_catalog_template_product_attributes_template FOREIGN KEY (product_template_id) REFERENCES catalog_product_templates(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_catalog_template_product_attributes_attribute FOREIGN KEY (product_attribute_id) REFERENCES catalog_product_attributes(id) ON DELETE RESTRICT;

ALTER TABLE catalog_product_price_histories
    ADD CONSTRAINT fk_catalog_product_price_histories_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT;

-- Cart module FKs
ALTER TABLE cart_cart_items
    ADD CONSTRAINT fk_cart_cart_items_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_cart_cart_items_vendor FOREIGN KEY (vendor_id) REFERENCES identity_vendors(id) ON DELETE RESTRICT;

-- Orders module FKs
ALTER TABLE orders_orders
    ADD CONSTRAINT fk_orders_orders_parent FOREIGN KEY (parent_id) REFERENCES orders_orders(id) ON DELETE RESTRICT;

ALTER TABLE orders_order_items
    ADD CONSTRAINT fk_orders_order_items_order FOREIGN KEY (order_id) REFERENCES orders_orders(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_orders_order_items_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT;

ALTER TABLE orders_order_histories
    ADD CONSTRAINT fk_orders_order_histories_order FOREIGN KEY (order_id) REFERENCES orders_orders(id) ON DELETE RESTRICT;

-- Payments module FKs
ALTER TABLE payments_payments
    ADD CONSTRAINT fk_payments_payments_order FOREIGN KEY (order_id) REFERENCES orders_orders(id) ON DELETE RESTRICT;

-- Shipping module FKs
ALTER TABLE shipping_shipments
    ADD CONSTRAINT fk_shipping_shipments_order FOREIGN KEY (order_id) REFERENCES orders_orders(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_shipping_shipments_warehouse FOREIGN KEY (warehouse_id) REFERENCES inventory_warehouses(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_shipping_shipments_vendor FOREIGN KEY (vendor_id) REFERENCES identity_vendors(id) ON DELETE RESTRICT;

ALTER TABLE shipping_shipment_items
    ADD CONSTRAINT fk_shipping_shipment_items_shipment FOREIGN KEY (shipment_id) REFERENCES shipping_shipments(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_shipping_shipment_items_order_item FOREIGN KEY (order_item_id) REFERENCES orders_order_items(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_shipping_shipment_items_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT;

-- Inventory module FKs
ALTER TABLE inventory_stocks
    ADD CONSTRAINT fk_inventory_stocks_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_inventory_stocks_warehouse FOREIGN KEY (warehouse_id) REFERENCES inventory_warehouses(id) ON DELETE RESTRICT;

ALTER TABLE inventory_stock_histories
    ADD CONSTRAINT fk_inventory_stock_histories_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_inventory_stock_histories_warehouse FOREIGN KEY (warehouse_id) REFERENCES inventory_warehouses(id) ON DELETE RESTRICT;

ALTER TABLE inventory_warehouses
    ADD CONSTRAINT fk_inventory_warehouses_address FOREIGN KEY (address_id) REFERENCES identity_addresses(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_inventory_warehouses_vendor FOREIGN KEY (vendor_id) REFERENCES identity_vendors(id) ON DELETE RESTRICT;

ALTER TABLE inventory_product_back_in_stock_subscriptions
    ADD CONSTRAINT fk_inventory_subscriptions_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT;

-- Pricing module FKs
ALTER TABLE pricing_coupons
    ADD CONSTRAINT fk_pricing_coupons_cart_rule FOREIGN KEY (cart_rule_id) REFERENCES pricing_cart_rules(id) ON DELETE RESTRICT;

ALTER TABLE pricing_cart_rule_categories
    ADD CONSTRAINT fk_pricing_cart_rule_categories_rule FOREIGN KEY (cart_rule_id) REFERENCES pricing_cart_rules(id) ON DELETE RESTRICT;

ALTER TABLE pricing_cart_rule_products
    ADD CONSTRAINT fk_pricing_cart_rule_products_rule FOREIGN KEY (cart_rule_id) REFERENCES pricing_cart_rules(id) ON DELETE RESTRICT;

ALTER TABLE pricing_cart_rule_customer_groups
    ADD CONSTRAINT fk_pricing_cart_rule_customer_groups_rule FOREIGN KEY (cart_rule_id) REFERENCES pricing_cart_rules(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_pricing_cart_rule_customer_groups_group FOREIGN KEY (customer_group_id) REFERENCES identity_customer_groups(id) ON DELETE RESTRICT;

ALTER TABLE pricing_cart_rule_usages
    ADD CONSTRAINT fk_pricing_cart_rule_usages_rule FOREIGN KEY (cart_rule_id) REFERENCES pricing_cart_rules(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_pricing_cart_rule_usages_coupon FOREIGN KEY (coupon_id) REFERENCES pricing_coupons(id) ON DELETE RESTRICT;

ALTER TABLE pricing_catalog_rule_customer_groups
    ADD CONSTRAINT fk_pricing_catalog_rule_customer_groups_rule FOREIGN KEY (catalog_rule_id) REFERENCES pricing_catalog_rules(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_pricing_catalog_rule_customer_groups_group FOREIGN KEY (customer_group_id) REFERENCES identity_customer_groups(id) ON DELETE RESTRICT;

-- Tax module FKs
ALTER TABLE tax_tax_rates
    ADD CONSTRAINT fk_tax_tax_rates_class FOREIGN KEY (tax_class_id) REFERENCES tax_tax_classes(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_tax_tax_rates_country FOREIGN KEY (country_id) REFERENCES identity_countries(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_tax_tax_rates_state_or_province FOREIGN KEY (state_or_province_id) REFERENCES identity_state_or_provinces(id) ON DELETE RESTRICT;

-- Reviews module FKs
ALTER TABLE reviews_reviews
    ADD CONSTRAINT fk_reviews_reviews_user FOREIGN KEY (user_id) REFERENCES identity_users(id) ON DELETE RESTRICT;

ALTER TABLE reviews_replies
    ADD CONSTRAINT fk_reviews_replies_review FOREIGN KEY (review_id) REFERENCES reviews_reviews(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_reviews_replies_user FOREIGN KEY (user_id) REFERENCES identity_users(id) ON DELETE RESTRICT;

-- CMS module FKs
ALTER TABLE cms_menu_items
    ADD CONSTRAINT fk_cms_menu_items_parent FOREIGN KEY (parent_id) REFERENCES cms_menu_items(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_cms_menu_items_menu FOREIGN KEY (menu_id) REFERENCES cms_menus(id) ON DELETE RESTRICT;

-- Activity log module FKs
ALTER TABLE activity_log_activities
    ADD CONSTRAINT fk_activity_log_activities_type FOREIGN KEY (activity_type_id) REFERENCES activity_log_activity_types(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_activity_log_activities_user FOREIGN KEY (user_id) REFERENCES identity_users(id) ON DELETE RESTRICT;

-- Notifications module FKs
ALTER TABLE notifications_notifications
    ADD CONSTRAINT fk_notifications_notifications_user FOREIGN KEY (user_id) REFERENCES identity_users(id) ON DELETE RESTRICT;

-- Comments module FKs
ALTER TABLE comments_comments
    ADD CONSTRAINT fk_comments_comments_user FOREIGN KEY (user_id) REFERENCES identity_users(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_comments_comments_parent FOREIGN KEY (parent_id) REFERENCES comments_comments(id) ON DELETE RESTRICT;

-- News module FKs
ALTER TABLE news_news_items
    ADD CONSTRAINT fk_news_news_items_thumbnail FOREIGN KEY (thumbnail_image_id) REFERENCES catalog_media(id) ON DELETE RESTRICT;

ALTER TABLE news_news_item_categories
    ADD CONSTRAINT fk_news_news_item_categories_item FOREIGN KEY (news_item_id) REFERENCES news_news_items(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_news_news_item_categories_category FOREIGN KEY (news_category_id) REFERENCES news_news_categories(id) ON DELETE RESTRICT;

-- Product comparison module FKs
ALTER TABLE product_comparison_comparing_products
    ADD CONSTRAINT fk_product_comparison_products_user FOREIGN KEY (user_id) REFERENCES identity_users(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_product_comparison_products_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT;

-- Wish list module FKs
ALTER TABLE wish_list_wish_lists
    ADD CONSTRAINT fk_wish_list_wish_lists_user FOREIGN KEY (user_id) REFERENCES identity_users(id) ON DELETE RESTRICT;

ALTER TABLE wish_list_wish_list_items
    ADD CONSTRAINT fk_wish_list_wish_list_items_list FOREIGN KEY (wish_list_id) REFERENCES wish_list_wish_lists(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_wish_list_wish_list_items_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT;

-- Recently viewed module FKs
ALTER TABLE product_recently_viewed_recently_viewed_products
    ADD CONSTRAINT fk_recently_viewed_products_user FOREIGN KEY (user_id) REFERENCES identity_users(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_recently_viewed_products_product FOREIGN KEY (product_id) REFERENCES catalog_products(id) ON DELETE RESTRICT;

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

-- Identity module indexes
CREATE INDEX IF NOT EXISTS idx_identity_users_email ON identity_users(email);
CREATE INDEX IF NOT EXISTS idx_identity_users_user_guid ON identity_users(user_guid);
CREATE INDEX IF NOT EXISTS idx_identity_users_deleted_at ON identity_users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_identity_users_vendor_id ON identity_users(vendor_id);
CREATE INDEX IF NOT EXISTS idx_identity_roles_name ON identity_roles(name);
CREATE INDEX IF NOT EXISTS idx_identity_user_roles_user_id ON identity_user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_identity_user_roles_role_id ON identity_user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_identity_addresses_state_or_province_id ON identity_addresses(state_or_province_id);
CREATE INDEX IF NOT EXISTS idx_identity_addresses_country_id ON identity_addresses(country_id);
CREATE INDEX IF NOT EXISTS idx_identity_addresses_district_id ON identity_addresses(district_id);
CREATE INDEX IF NOT EXISTS idx_identity_user_addresses_user_id ON identity_user_addresses(user_id);
CREATE INDEX IF NOT EXISTS idx_identity_user_addresses_address_id ON identity_user_addresses(address_id);
CREATE INDEX IF NOT EXISTS idx_identity_state_or_provinces_country_id ON identity_state_or_provinces(country_id);
CREATE INDEX IF NOT EXISTS idx_identity_districts_state_or_province_id ON identity_districts(state_or_province_id);
CREATE INDEX IF NOT EXISTS idx_identity_vendors_slug ON identity_vendors(slug);
CREATE INDEX IF NOT EXISTS idx_identity_customer_groups_name ON identity_customer_groups(name);
CREATE INDEX IF NOT EXISTS idx_identity_customer_group_users_user_id ON identity_customer_group_users(user_id);
CREATE INDEX IF NOT EXISTS idx_identity_customer_group_users_group_id ON identity_customer_group_users(customer_group_id);
CREATE INDEX IF NOT EXISTS idx_identity_media_deleted_at ON identity_media(deleted_at);
CREATE INDEX IF NOT EXISTS idx_identity_widget_instances_widget_id ON identity_widget_instances(widget_id);
CREATE INDEX IF NOT EXISTS idx_identity_widget_instances_widget_zone_id ON identity_widget_instances(widget_zone_id);

-- Catalog module indexes
CREATE INDEX IF NOT EXISTS idx_catalog_products_slug ON catalog_products(slug);
CREATE INDEX IF NOT EXISTS idx_catalog_products_brand_id ON catalog_products(brand_id);
CREATE INDEX IF NOT EXISTS idx_catalog_products_vendor_id ON catalog_products(vendor_id);
CREATE INDEX IF NOT EXISTS idx_catalog_products_tax_class_id ON catalog_products(tax_class_id);
CREATE INDEX IF NOT EXISTS idx_catalog_products_thumbnail_image_id ON catalog_products(thumbnail_image_id);
CREATE INDEX IF NOT EXISTS idx_catalog_products_is_published ON catalog_products(is_published);
CREATE INDEX IF NOT EXISTS idx_catalog_products_is_featured ON catalog_products(is_featured);
CREATE INDEX IF NOT EXISTS idx_catalog_products_deleted_at ON catalog_products(deleted_at);
CREATE INDEX IF NOT EXISTS idx_catalog_products_sku ON catalog_products(sku);
CREATE INDEX IF NOT EXISTS idx_catalog_categories_slug ON catalog_categories(slug);
CREATE INDEX IF NOT EXISTS idx_catalog_categories_parent_id ON catalog_categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_catalog_categories_thumbnail_image_id ON catalog_categories(thumbnail_image_id);
CREATE INDEX IF NOT EXISTS idx_catalog_categories_deleted_at ON catalog_categories(deleted_at);
CREATE INDEX IF NOT EXISTS idx_catalog_brands_slug ON catalog_brands(slug);
CREATE INDEX IF NOT EXISTS idx_catalog_product_categories_product_id ON catalog_product_categories(product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_categories_category_id ON catalog_product_categories(category_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_attributes_group_id ON catalog_product_attributes(group_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_attribute_values_attribute_id ON catalog_product_attribute_values(attribute_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_attribute_values_product_id ON catalog_product_attribute_values(product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_option_values_option_id ON catalog_product_option_values(option_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_option_values_product_id ON catalog_product_option_values(product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_option_combinations_product_id ON catalog_product_option_combinations(product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_option_combinations_option_id ON catalog_product_option_combinations(option_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_links_product_id ON catalog_product_links(product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_links_linked_product_id ON catalog_product_links(linked_product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_medias_product_id ON catalog_product_medias(product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_medias_media_id ON catalog_product_medias(media_id);
CREATE INDEX IF NOT EXISTS idx_catalog_product_price_histories_product_id ON catalog_product_price_histories(product_id);
CREATE INDEX IF NOT EXISTS idx_catalog_media_deleted_at ON catalog_media(deleted_at);

-- Cart module indexes
CREATE INDEX IF NOT EXISTS idx_cart_cart_items_customer_id ON cart_cart_items(customer_id);
CREATE INDEX IF NOT EXISTS idx_cart_cart_items_product_id ON cart_cart_items(product_id);
CREATE INDEX IF NOT EXISTS idx_cart_cart_items_vendor_id ON cart_cart_items(vendor_id);

-- Orders module indexes
CREATE INDEX IF NOT EXISTS idx_orders_orders_customer_id ON orders_orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_orders_order_status ON orders_orders(order_status);
CREATE INDEX IF NOT EXISTS idx_orders_orders_parent_id ON orders_orders(parent_id);
CREATE INDEX IF NOT EXISTS idx_orders_orders_created_at ON orders_orders(created_at);
CREATE INDEX IF NOT EXISTS idx_orders_order_items_order_id ON orders_order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_orders_order_items_product_id ON orders_order_items(product_id);
CREATE INDEX IF NOT EXISTS idx_orders_order_histories_order_id ON orders_order_histories(order_id);

-- Payments module indexes
CREATE INDEX IF NOT EXISTS idx_payments_payments_order_id ON payments_payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_payments_status ON payments_payments(status);

-- Shipping module indexes
CREATE INDEX IF NOT EXISTS idx_shipping_shipments_order_id ON shipping_shipments(order_id);
CREATE INDEX IF NOT EXISTS idx_shipping_shipments_warehouse_id ON shipping_shipments(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_shipping_shipments_vendor_id ON shipping_shipments(vendor_id);
CREATE INDEX IF NOT EXISTS idx_shipping_shipment_items_shipment_id ON shipping_shipment_items(shipment_id);
CREATE INDEX IF NOT EXISTS idx_shipping_shipment_items_order_item_id ON shipping_shipment_items(order_item_id);
CREATE INDEX IF NOT EXISTS idx_shipping_shipment_items_product_id ON shipping_shipment_items(product_id);

-- Inventory module indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_inventory_stocks_product_warehouse ON inventory_stocks(product_id, warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inventory_stocks_product_id ON inventory_stocks(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_stocks_warehouse_id ON inventory_stocks(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inventory_stock_histories_product_id ON inventory_stock_histories(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_stock_histories_warehouse_id ON inventory_stock_histories(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inventory_warehouses_address_id ON inventory_warehouses(address_id);
CREATE INDEX IF NOT EXISTS idx_inventory_warehouses_vendor_id ON inventory_warehouses(vendor_id);
CREATE INDEX IF NOT EXISTS idx_inventory_subscriptions_product_id ON inventory_product_back_in_stock_subscriptions(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_subscriptions_email ON inventory_product_back_in_stock_subscriptions(customer_email);

-- Pricing module indexes
CREATE INDEX IF NOT EXISTS idx_pricing_coupons_cart_rule_id ON pricing_coupons(cart_rule_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_pricing_coupons_code ON pricing_coupons(code);
CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_categories_rule_id ON pricing_cart_rule_categories(cart_rule_id);
CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_products_rule_id ON pricing_cart_rule_products(cart_rule_id);
CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_customer_groups_rule_id ON pricing_cart_rule_customer_groups(cart_rule_id);
CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_customer_groups_group_id ON pricing_cart_rule_customer_groups(customer_group_id);
CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_usages_rule_id ON pricing_cart_rule_usages(cart_rule_id);
CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_usages_coupon_id ON pricing_cart_rule_usages(coupon_id);
CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_usages_order_id ON pricing_cart_rule_usages(order_id);
CREATE INDEX IF NOT EXISTS idx_pricing_cart_rule_usages_user_id ON pricing_cart_rule_usages(user_id);
CREATE INDEX IF NOT EXISTS idx_pricing_catalog_rule_customer_groups_rule_id ON pricing_catalog_rule_customer_groups(catalog_rule_id);
CREATE INDEX IF NOT EXISTS idx_pricing_catalog_rule_customer_groups_group_id ON pricing_catalog_rule_customer_groups(customer_group_id);

-- Tax module indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_tax_tax_classes_name ON tax_tax_classes(name);
CREATE INDEX IF NOT EXISTS idx_tax_tax_rates_tax_class_id ON tax_tax_rates(tax_class_id);
CREATE INDEX IF NOT EXISTS idx_tax_tax_rates_country_id ON tax_tax_rates(country_id);
CREATE INDEX IF NOT EXISTS idx_tax_tax_rates_state_or_province_id ON tax_tax_rates(state_or_province_id);

-- Reviews module indexes
CREATE INDEX IF NOT EXISTS idx_reviews_reviews_user_id ON reviews_reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_reviews_entity_type_id ON reviews_reviews(entity_type_id);
CREATE INDEX IF NOT EXISTS idx_reviews_reviews_entity_id ON reviews_reviews(entity_id);
CREATE INDEX IF NOT EXISTS idx_reviews_reviews_status ON reviews_reviews(status);
CREATE INDEX IF NOT EXISTS idx_reviews_replies_review_id ON reviews_replies(review_id);
CREATE INDEX IF NOT EXISTS idx_reviews_replies_user_id ON reviews_replies(user_id);

-- CMS module indexes
CREATE INDEX IF NOT EXISTS idx_cms_pages_slug ON cms_pages(slug);
CREATE INDEX IF NOT EXISTS idx_cms_menu_items_menu_id ON cms_menu_items(menu_id);
CREATE INDEX IF NOT EXISTS idx_cms_menu_items_parent_id ON cms_menu_items(parent_id);
CREATE INDEX IF NOT EXISTS idx_cms_menu_items_entity_id ON cms_menu_items(entity_id);

-- Activity log module indexes
CREATE INDEX IF NOT EXISTS idx_activity_log_activities_type_id ON activity_log_activities(activity_type_id);
CREATE INDEX IF NOT EXISTS idx_activity_log_activities_user_id ON activity_log_activities(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_log_activities_entity_type ON activity_log_activities(entity_type_id);

-- Notifications module indexes
CREATE INDEX IF NOT EXISTS idx_notifications_notifications_user_id ON notifications_notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_notifications_is_read ON notifications_notifications(is_read);

-- Search module indexes
CREATE INDEX IF NOT EXISTS idx_search_queries_query_text ON search_queries(query_text);
CREATE INDEX IF NOT EXISTS idx_search_queries_created_at ON search_queries(created_at);

-- Comments module indexes
CREATE INDEX IF NOT EXISTS idx_comments_comments_user_id ON comments_comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_comments_parent_id ON comments_comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_comments_entity_type_id ON comments_comments(entity_type_id);
CREATE INDEX IF NOT EXISTS idx_comments_comments_entity_id ON comments_comments(entity_id);

-- News module indexes
CREATE INDEX IF NOT EXISTS idx_news_news_items_slug ON news_news_items(slug);
CREATE INDEX IF NOT EXISTS idx_news_news_items_thumbnail_image_id ON news_news_items(thumbnail_image_id);
CREATE INDEX IF NOT EXISTS idx_news_news_item_categories_item_id ON news_news_item_categories(news_item_id);
CREATE INDEX IF NOT EXISTS idx_news_news_item_categories_category_id ON news_news_item_categories(news_category_id);

-- Product comparison indexes
CREATE INDEX IF NOT EXISTS idx_product_comparison_products_user_id ON product_comparison_comparing_products(user_id);
CREATE INDEX IF NOT EXISTS idx_product_comparison_products_product_id ON product_comparison_comparing_products(product_id);

-- Wish list indexes
CREATE INDEX IF NOT EXISTS idx_wish_list_wish_lists_user_id ON wish_list_wish_lists(user_id);
CREATE INDEX IF NOT EXISTS idx_wish_list_wish_list_items_wish_list_id ON wish_list_wish_list_items(wish_list_id);
CREATE INDEX IF NOT EXISTS idx_wish_list_wish_list_items_product_id ON wish_list_wish_list_items(product_id);

-- Recently viewed indexes
CREATE INDEX IF NOT EXISTS idx_recently_viewed_products_user_id ON product_recently_viewed_recently_viewed_products(user_id);
CREATE INDEX IF NOT EXISTS idx_recently_viewed_products_product_id ON product_recently_viewed_recently_viewed_products(product_id);
CREATE INDEX IF NOT EXISTS idx_recently_viewed_products_viewed_on ON product_recently_viewed_recently_viewed_products(latest_viewed_on);
