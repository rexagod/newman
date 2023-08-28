CREATE TABLE IF NOT EXISTS snipes (
    message_id BIGINT PRIMARY KEY,
    channel_id BIGINT NOT NULL,
    author_id BIGINT NOT NULL,
    content JSONB DEFAULT '{}'::jsonb,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
