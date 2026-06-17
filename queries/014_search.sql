-- ============================================================
-- SEARCH QUERIES
-- ============================================================
CREATE TABLE IF NOT EXISTS search_queries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query VARCHAR(500) NOT NULL,
    result_count INT DEFAULT 0 CHECK (result_count >= 0),
    user_id UUID REFERENCES identity_users(id) ON DELETE RESTRICT,
    ip_address VARCHAR(45),
    source VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE search_queries IS 'Logged search queries for analytics and suggestions';

CREATE INDEX IF NOT EXISTS idx_search_queries_user_id ON search_queries(user_id);
CREATE INDEX IF NOT EXISTS idx_search_queries_query ON search_queries(query);
