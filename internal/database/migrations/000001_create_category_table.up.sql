CREATE TABLE IF NOT EXISTS category (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id SERIAL,
    name VARCHAR(255) UNIQUE
);