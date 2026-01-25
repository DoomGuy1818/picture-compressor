DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_type
        WHERE typname = 'picture_status'
    ) THEN
CREATE TYPE picture_status AS ENUM (
            'pending',
            'processing',
            'uploaded',
            'failed'
        );
END IF;
END$$;

CREATE TABLE IF NOT EXISTS pictures (
    id UUID PRIMARY KEY,
    file_path TEXT,
    status picture_status NOT NULL DEFAULT 'pending'
);