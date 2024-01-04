-- Drop the foreign key constraints in 'post_tag' table first
ALTER TABLE post_tag DROP CONSTRAINT post_tag_post_id_fkey;

-- Drop the 'post_tag' table if it exists.
DROP TABLE IF EXISTS post_tag;
