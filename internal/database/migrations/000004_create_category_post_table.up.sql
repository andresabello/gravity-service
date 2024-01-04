CREATE TABLE IF NOT EXISTS category_post (
    post_id UUID,
    category_id UUID,
    FOREIGN KEY (post_id) REFERENCES post (id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, category_id)
);