CREATE TABLE IF NOT EXISTS sources (
    source_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    source_type VARCHAR(255)
    -- Add other columns as needed
);
