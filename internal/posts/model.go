package posts

import (
	"database/sql"
	"pi-search/pkg/timeutil"
	"pi-search/pkg/tracer"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Post represents a blog post from Wordpress
type Post struct {
	ID            uuid.UUID           `gorm:"id"`
	PostID        uint                `json:"id" gorm:"column:post_id"`
	Date          timeutil.CustomTime `json:"date" gorm:"column:date"`
	Slug          string              `json:"slug" gorm:"slug"`
	Status        string              `json:"status" gorm:"status"`
	Type          string              `json:"type" gorm:"type"`
	Link          string              `json:"link" gorm:"link"`
	Title         string              `json:"title" gorm:"column:title"`
	Content       string              `json:"content" gorm:"column:content"`
	Excerpt       string              `json:"excerpt" gorm:"column:excerpt"`
	Author        uint                `json:"author" gorm:"column:author_name"`
	FeaturedMedia *int                `json:"featured_media" gorm:"featured_media"`
	// TODO: add categories and tags
	// Categories    []uint     `json:"categories" db:"categories"`
	// Tags          []uint     `json:"tags" db:"tags"`
	Source string `json:"source" gorm:"source"`
}

// CreatePost inserts a new post into the database.
func CreatePost(db *gorm.DB, post *Post) error {
	result := db.Create(post)

	if result.Error != nil {
		return tracer.TraceError(result.Error)
	}

	return nil
}

// GetPost retrieves a post from the database by ID.
func GetPost(db *gorm.DB, id uint, source string) (Post, error) {
	var post Post
	db.Where("post_id=? AND source=?", id, source).First(&post)
	return post, nil
}

// GetPostsByQueryString retrieves all posts that contain queryString in
// Content and Title.
func GetPostsByQueryString(db *gorm.DB, queryString string, page string) ([]Post, error) {
	var posts []Post
	db.Raw(
		`SELECT * 
		FROM post,
		to_tsvector(post.title || post.excerpt) document,
		to_tsquery(@query) query,
		NULLIF(ts_rank(to_tsvector(post.title), query), 0) rank_title,
    	NULLIF(ts_rank(to_tsvector(post.excerpt), query), 0) rank_excerpt,
		SIMILARITY(@query, post.title || post.excerpt) similarity
		WHERE query @@ document OR similarity > 0
		ORDER BY rank_title, rank_excerpt, similarity DESC NULLS LAST
		LIMIT 10 OFFSET @page;`,
		sql.Named("query", queryString),
		sql.Named("page", page),
	).Scan(&posts)

	return posts, nil
}

// GetAllPosts retrieves all posts from the database.
func GetAllPosts(db *gorm.DB) ([]Post, error) {
	var posts []Post
	db.Find(&posts)
	return posts, nil
}

// UpdatePost updates a post in the database.
func UpdatePost(db *gorm.DB, post *Post) error {
	db.Save(post)
	// return tracer.TraceError(err)
	return nil
}

// DeletePost deletes a post from the database by UUID.
func DeletePost(db *gorm.DB, id uint) error {
	db.Delete(&Post{}, id)
	return nil
}
