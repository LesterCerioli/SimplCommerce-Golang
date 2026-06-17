-- ============================================================
-- PAGES
-- ============================================================
CREATE TABLE IF NOT EXISTS cms_pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    slug VARCHAR(500) NOT NULL UNIQUE,
    content TEXT,
    meta_title VARCHAR(255),
    meta_description TEXT,
    is_active BOOLEAN DEFAULT true,
    is_published BOOLEAN DEFAULT false,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE cms_pages IS 'CMS content pages';

CREATE INDEX IF NOT EXISTS idx_cms_pages_slug ON cms_pages(slug);

-- ============================================================
-- MENUS
-- ============================================================
CREATE TABLE IF NOT EXISTS cms_menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE cms_menus IS 'Navigation menus';

-- ============================================================
-- MENU ITEMS
-- ============================================================
CREATE TABLE IF NOT EXISTS cms_menu_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    menu_id UUID NOT NULL REFERENCES cms_menus(id) ON DELETE RESTRICT,
    parent_id UUID REFERENCES cms_menu_items(id) ON DELETE RESTRICT,
    entity_id UUID,
    title VARCHAR(500) NOT NULL,
    url TEXT,
    target VARCHAR(50) DEFAULT '_self',
    icon VARCHAR(255),
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE cms_menu_items IS 'Items within a navigation menu';

CREATE INDEX IF NOT EXISTS idx_cms_menu_items_menu_id ON cms_menu_items(menu_id);
CREATE INDEX IF NOT EXISTS idx_cms_menu_items_parent_id ON cms_menu_items(parent_id);
CREATE INDEX IF NOT EXISTS idx_cms_menu_items_entity_id ON cms_menu_items(entity_id);
