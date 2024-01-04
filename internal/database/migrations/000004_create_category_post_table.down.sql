-- Drop the foreign key constraints in 'category_post' table first
ALTER TABLE category_post DROP CONSTRAINT category_post_post_id_fkey;
ALTER TABLE category_post DROP CONSTRAINT category_post_category_id_fkey;

-- Drop the 'category_post' table if it exists.
DROP TABLE IF EXISTS category_post;
