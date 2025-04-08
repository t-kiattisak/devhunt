CREATE TABLE tool_votes (
    id SERIAL PRIMARY KEY,
    tool_id INTEGER NOT NULL REFERENCES tools(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(tool_id, user_id)
);