-- Create models table
CREATE TABLE Models (
    id SERIAL PRIMARY KEY,
    make_id INT,
    model_name VARCHAR(255) NOT NULL UNIQUE,
    FOREIGN KEY (make_id) REFERENCES Makes(id)
);