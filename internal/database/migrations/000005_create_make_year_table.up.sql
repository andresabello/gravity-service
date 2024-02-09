-- Create make_years pivot table
CREATE TABLE make_years (
    make_id INT,
    year_id INT,
    PRIMARY KEY (make_id, year_id),
    FOREIGN KEY (make_id) REFERENCES Makes(id),
    FOREIGN KEY (year_id) REFERENCES Years(id)
);