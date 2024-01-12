CREATE TABLE IF NOT EXISTS categories_post (
    post_id UUID,
    category_id UUID,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, category_id)
);