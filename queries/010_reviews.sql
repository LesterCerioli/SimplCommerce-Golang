-- ============================================================
-- REVIEWS
-- ============================================================
CREATE TABLE IF NOT EXISTS reviews_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
    entity_type_id UUID NOT NULL,
    entity_id UUID NOT NULL,
    title VARCHAR(500),
    content TEXT NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    is_approved BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE reviews_reviews IS 'Product/service reviews from users';

CREATE INDEX IF NOT EXISTS idx_reviews_reviews_user_id ON reviews_reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_reviews_entity_type_id ON reviews_reviews(entity_type_id);
CREATE INDEX IF NOT EXISTS idx_reviews_reviews_entity_id ON reviews_reviews(entity_id);

-- ============================================================
-- REVIEW REPLIES
-- ============================================================
CREATE TABLE IF NOT EXISTS reviews_replies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id UUID NOT NULL REFERENCES reviews_reviews(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
    content TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE reviews_replies IS 'Replies to reviews';

CREATE INDEX IF NOT EXISTS idx_reviews_replies_review_id ON reviews_replies(review_id);
CREATE INDEX IF NOT EXISTS idx_reviews_replies_user_id ON reviews_replies(user_id);
