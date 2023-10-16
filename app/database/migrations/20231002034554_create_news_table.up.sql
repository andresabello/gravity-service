CREATE TABLE IF NOT EXISTS news (
    news_id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    content TEXT,
    publish_date TIMESTAMP,
    source_id INT, -- Replace with actual foreign key reference to sources table
    team_id INT -- Replace with actual foreign key reference to teams table
    -- Add other columns as needed
);
