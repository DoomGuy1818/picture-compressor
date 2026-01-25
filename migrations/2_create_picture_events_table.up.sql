CREATE TABLE IF NOT EXISTS pictures_events (
    id UUID PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    payload TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'new' CHECK(status in ('new', 'done')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reserved_to TIMESTAMP DEFAULT NULL
);