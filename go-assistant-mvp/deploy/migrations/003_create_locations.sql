CREATE TABLE IF NOT EXISTS location_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    lat DOUBLE PRECISION,
    lng DOUBLE PRECISION,
    recorded_at TIMESTAMPTZ DEFAULT NOW()
);
