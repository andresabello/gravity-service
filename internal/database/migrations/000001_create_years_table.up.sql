-- Create years table
CREATE TABLE years (
    id SERIAL PRIMARY KEY,
    year INT NOT NULL UNIQUE
);