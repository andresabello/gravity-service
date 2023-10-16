CREATE TABLE IF NOT EXISTS teams (
    team_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    handle VARCHAR(255),
    sport_id INT
    -- Replace with actual foreign key reference to sports table
    -- Add other columns as needed
);
