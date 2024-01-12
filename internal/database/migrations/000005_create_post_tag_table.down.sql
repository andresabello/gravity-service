-- Drop the foreign key constraints in 'post_tag' table first
ALTER TABLE posts_tag DROP CONSTRAINT posts_tag_post_id_fkey;

-- Drop the 'post_tag' table if it exists.
DROP TABLE IF EXISTS posts_tag;
