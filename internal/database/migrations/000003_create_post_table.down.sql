-- Drop trigram indexes for each column
DROP INDEX IF EXISTS idx_title_search;
DROP INDEX IF EXISTS idx_excerpt_search;
DROP INDEX IF EXISTS idx_content_search;


-- Drop the 'post' table if it exists.
DROP TABLE IF EXISTS post;
