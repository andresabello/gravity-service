CREATE TABLE IF NOT EXISTS posts_tag (
    post_id UUID,
    tag_id UUID,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);