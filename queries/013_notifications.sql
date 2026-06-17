-- ============================================================
-- NOTIFICATIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS notifications_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
    title VARCHAR(500) NOT NULL,
    body TEXT,
    type VARCHAR(100),
    reference_type VARCHAR(255),
    reference_id UUID,
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE notifications_notifications IS 'User notifications';

CREATE INDEX IF NOT EXISTS idx_notifications_notifications_user_id ON notifications_notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_notifications_is_read ON notifications_notifications(is_read);
