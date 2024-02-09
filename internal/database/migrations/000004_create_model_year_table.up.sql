-- Create model_years pivot table
CREATE TABLE model_years (
    model_id INT,
    year_id INT,
    PRIMARY KEY (model_id, year_id),
    FOREIGN KEY (model_id) REFERENCES Models(id),
    FOREIGN KEY (year_id) REFERENCES Years(id)
);