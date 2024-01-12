CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tag_id SERIAL,
    name VARCHAR(255) UNIQUE
);