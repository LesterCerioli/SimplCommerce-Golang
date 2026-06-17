-- ============================================================
-- TRIGGER: auto_update_updated_at
-- Automatically updates updated_at on row modification
-- ============================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Identity tables
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_identity_countries_updated_at') THEN
CREATE TRIGGER trg_identity_countries_updated_at BEFORE UPDATE ON identity_countries FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_identity_roles_updated_at') THEN
CREATE TRIGGER trg_identity_roles_updated_at BEFORE UPDATE ON identity_roles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_identity_users_updated_at') THEN
CREATE TRIGGER trg_identity_users_updated_at BEFORE UPDATE ON identity_users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_identity_addresses_updated_at') THEN
CREATE TRIGGER trg_identity_addresses_updated_at BEFORE UPDATE ON identity_addresses FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_identity_state_or_provinces_updated_at') THEN
CREATE TRIGGER trg_identity_state_or_provinces_updated_at BEFORE UPDATE ON identity_state_or_provinces FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_identity_districts_updated_at') THEN
CREATE TRIGGER trg_identity_districts_updated_at BEFORE UPDATE ON identity_districts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_identity_vendors_updated_at') THEN
CREATE TRIGGER trg_identity_vendors_updated_at BEFORE UPDATE ON identity_vendors FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_identity_customer_groups_updated_at') THEN
CREATE TRIGGER trg_identity_customer_groups_updated_at BEFORE UPDATE ON identity_customer_groups FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_identity_media_updated_at') THEN
CREATE TRIGGER trg_identity_media_updated_at BEFORE UPDATE ON identity_media FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;

-- Catalog tables
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_catalog_categories_updated_at') THEN
CREATE TRIGGER trg_catalog_categories_updated_at BEFORE UPDATE ON catalog_categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_catalog_brands_updated_at') THEN
CREATE TRIGGER trg_catalog_brands_updated_at BEFORE UPDATE ON catalog_brands FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_catalog_products_updated_at') THEN
CREATE TRIGGER trg_catalog_products_updated_at BEFORE UPDATE ON catalog_products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_catalog_media_updated_at') THEN
CREATE TRIGGER trg_catalog_media_updated_at BEFORE UPDATE ON catalog_media FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_catalog_product_attribute_groups_updated_at') THEN
CREATE TRIGGER trg_catalog_product_attribute_groups_updated_at BEFORE UPDATE ON catalog_product_attribute_groups FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_catalog_product_attributes_updated_at') THEN
CREATE TRIGGER trg_catalog_product_attributes_updated_at BEFORE UPDATE ON catalog_product_attributes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_catalog_product_options_updated_at') THEN
CREATE TRIGGER trg_catalog_product_options_updated_at BEFORE UPDATE ON catalog_product_options FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_catalog_product_templates_updated_at') THEN
CREATE TRIGGER trg_catalog_product_templates_updated_at BEFORE UPDATE ON catalog_product_templates FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;

-- Cart
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_cart_cart_items_updated_at') THEN
CREATE TRIGGER trg_cart_cart_items_updated_at BEFORE UPDATE ON cart_cart_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;

-- Orders
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_orders_orders_updated_at') THEN
CREATE TRIGGER trg_orders_orders_updated_at BEFORE UPDATE ON orders_orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_orders_order_addresses_updated_at') THEN
CREATE TRIGGER trg_orders_order_addresses_updated_at BEFORE UPDATE ON orders_order_addresses FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;

-- Payments
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_payments_payments_updated_at') THEN
CREATE TRIGGER trg_payments_payments_updated_at BEFORE UPDATE ON payments_payments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;

-- Shipping
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_shipping_shipments_updated_at') THEN
CREATE TRIGGER trg_shipping_shipments_updated_at BEFORE UPDATE ON shipping_shipments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_shipping_shipping_providers_updated_at') THEN
CREATE TRIGGER trg_shipping_shipping_providers_updated_at BEFORE UPDATE ON shipping_shipping_providers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_shipping_table_rate_price_destinations_updated_at') THEN
CREATE TRIGGER trg_shipping_table_rate_price_destinations_updated_at BEFORE UPDATE ON shipping_table_rate_price_and_destinations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;

-- Inventory
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_inventory_warehouses_updated_at') THEN
CREATE TRIGGER trg_inventory_warehouses_updated_at BEFORE UPDATE ON inventory_warehouses FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_inventory_stocks_updated_at') THEN
CREATE TRIGGER trg_inventory_stocks_updated_at BEFORE UPDATE ON inventory_stocks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;

-- Pricing
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_pricing_cart_rules_updated_at') THEN
CREATE TRIGGER trg_pricing_cart_rules_updated_at BEFORE UPDATE ON pricing_cart_rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_pricing_catalog_rules_updated_at') THEN
CREATE TRIGGER trg_pricing_catalog_rules_updated_at BEFORE UPDATE ON pricing_catalog_rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;

-- Tax
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_tax_tax_classes_updated_at') THEN
CREATE TRIGGER trg_tax_tax_classes_updated_at BEFORE UPDATE ON tax_tax_classes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_tax_tax_rates_updated_at') THEN
CREATE TRIGGER trg_tax_tax_rates_updated_at BEFORE UPDATE ON tax_tax_rates FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;

-- Reviews
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_reviews_reviews_updated_at') THEN
CREATE TRIGGER trg_reviews_reviews_updated_at BEFORE UPDATE ON reviews_reviews FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_reviews_replies_updated_at') THEN
CREATE TRIGGER trg_reviews_replies_updated_at BEFORE UPDATE ON reviews_replies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;

-- CMS
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_cms_pages_updated_at') THEN
CREATE TRIGGER trg_cms_pages_updated_at BEFORE UPDATE ON cms_pages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_cms_menus_updated_at') THEN
CREATE TRIGGER trg_cms_menus_updated_at BEFORE UPDATE ON cms_menus FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_cms_menu_items_updated_at') THEN
CREATE TRIGGER trg_cms_menu_items_updated_at BEFORE UPDATE ON cms_menu_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;

-- Notifications
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_notifications_notifications_updated_at') THEN
CREATE TRIGGER trg_notifications_notifications_updated_at BEFORE UPDATE ON notifications_notifications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); END IF; END $$;
