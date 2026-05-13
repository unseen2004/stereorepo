CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    due_at TIMESTAMPTZ,
    location_trigger JSONB,
    source VARCHAR(50),
    external_id TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
