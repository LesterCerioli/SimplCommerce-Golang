-- ============================================================
-- ACTIVITY TYPES
-- ============================================================
CREATE TABLE IF NOT EXISTS activity_log_activity_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE activity_log_activity_types IS 'Types of activities that can be logged';

-- ============================================================
-- ACTIVITIES
-- ============================================================
CREATE TABLE IF NOT EXISTS activity_log_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    activity_type_id UUID NOT NULL REFERENCES activity_log_activity_types(id) ON DELETE RESTRICT,
    user_id UUID REFERENCES identity_users(id) ON DELETE RESTRICT,
    subject_type VARCHAR(255),
    subject_id UUID,
    description TEXT NOT NULL,
    metadata JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE activity_log_activities IS 'Activity log entries recording user actions';

CREATE INDEX IF NOT EXISTS idx_activity_log_activities_activity_type_id ON activity_log_activities(activity_type_id);
CREATE INDEX IF NOT EXISTS idx_activity_log_activities_user_id ON activity_log_activities(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_log_activities_subject ON activity_log_activities(subject_type, subject_id);
