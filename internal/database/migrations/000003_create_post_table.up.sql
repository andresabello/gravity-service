CREATE TABLE IF NOT EXISTS posts (
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
    author_name INT NOT NULL,
    featured_media INT,
    source VARCHAR(255) NOT NULL
);

CREATE INDEX idx_title_search ON posts USING gin(title gin_trgm_ops);
CREATE INDEX idx_excerpt_search ON posts USING gin(excerpt gin_trgm_ops);
CREATE INDEX idx_content_search ON posts USING gin(content gin_trgm_ops);