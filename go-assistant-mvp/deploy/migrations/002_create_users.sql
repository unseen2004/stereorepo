CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    google_token JSONB,
    notion_token TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
