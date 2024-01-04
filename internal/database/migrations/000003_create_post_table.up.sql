-- Enable the uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable pg_trgm includes support for trigram indexing. 
-- Full-Text Search Extension.
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS post (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id SERIAL NOT NULL,
    date TIMESTAMP NOT NULL,
    slug VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    link VARCHAR(255) NOT NULL,
    title TEXT,
    content TEXT,
    excerpt TEXT,
    author_name VARCHAR(255) NOT NULL,
    featured_media INT,
    source VARCHAR(255) NOT NULL
);

CREATE INDEX idx_title_search ON post USING gin(title gin_trgm_ops);
CREATE INDEX idx_excerpt_search ON post USING gin(excerpt gin_trgm_ops);
CREATE INDEX idx_content_search ON post USING gin(content gin_trgm_ops);