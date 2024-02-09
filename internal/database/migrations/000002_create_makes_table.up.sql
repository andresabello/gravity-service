-- Create makes table
CREATE TABLE makes (
    id SERIAL PRIMARY KEY,
    make_name VARCHAR(255) NOT NULL UNIQUE
);