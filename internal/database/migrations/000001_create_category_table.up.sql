-- Enable the uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable pg_trgm includes support for trigram indexing. 
-- Full-Text Search Extension.
CREATE EXTENSION IF NOT EXISTS pg_trgm;


CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id SERIAL,
    name VARCHAR(255) UNIQUE
);