-- Drop the foreign key constraints in 'categories_post' table first
ALTER TABLE categories_post DROP CONSTRAINT categories_post_post_id_fkey;
ALTER TABLE categories_post DROP CONSTRAINT categories_post_category_id_fkey;

-- Drop the 'categories_post' table if it exists.
DROP TABLE IF EXISTS categories_post;
