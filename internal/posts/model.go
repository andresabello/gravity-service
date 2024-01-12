package posts

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"pi-search/pkg/timeutil"
	"pi-search/pkg/tracer"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Rendered struct {
	Rendered string `json:"rendered"`
}

// Post represents a blog post from Wordpress
type Post struct {
	ID            uuid.UUID           `gorm:"id"`
	PostID        uint                `json:"id" gorm:"column:post_id"`
	Date          timeutil.CustomTime `json:"date" gorm:"column:date"`
	Slug          string              `json:"slug" gorm:"slug"`
	Status        string              `json:"status" gorm:"status"`
	Type          string              `json:"type" gorm:"type"`
	Link          string              `json:"link" gorm:"link"`
	Title         sql.NullString      `json:"title" gorm:"column:title;default:null"`
	Content       sql.NullString      `json:"content" gorm:"column:content;default:null"`
	Excerpt       sql.NullString      `json:"excerpt" gorm:"column:excerpt;default:null"`
	Author        int                 `json:"author" gorm:"column:author_name;default:1"`
	FeaturedMedia *int                `json:"featured_media" gorm:"featured_media"`
	// TODO: add categories and tags
	// Categories    []uint     `json:"categories" db:"categories"`
	// Tags          []uint     `json:"tags" db:"tags"`
	Source string `json:"source" gorm:"source"`
}

func (p *Post) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Title   Rendered `json:"title" gorm:"column:title"`
		Content Rendered `json:"content" gorm:"column:content"`
		Excerpt Rendered `json:"excerpt" gorm:"column:excerpt"`

		// Add other fields here
		ID            uuid.UUID           `gorm:"id"`
		PostID        uint                `json:"id" gorm:"column:post_id"`
		Date          timeutil.CustomTime `json:"date" gorm:"column:date"`
		Slug          string              `json:"slug" gorm:"slug"`
		Status        string              `json:"status" gorm:"status"`
		Type          string              `json:"type" gorm:"type"`
		Link          string              `json:"link" gorm:"link"`
		Author        int                 `json:"author" gorm:"column:author_name"`
		FeaturedMedia *int                `json:"featured_media" gorm:"featured_media"`
		// TODO: add categories and tags
		// Categories    []uint     `json:"categories" db:"categories"`
		// Tags          []uint     `json:"tags" db:"tags"`
		Source string `json:"source" gorm:"source"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	p.Title = sql.NullString{String: "", Valid: false}
	if len(tmp.Title.Rendered) > 0 {
		p.Title = sql.NullString{String: tmp.Title.Rendered, Valid: true}
	}

	p.Content = sql.NullString{String: "", Valid: false}
	if len(tmp.Content.Rendered) > 0 {
		p.Content = sql.NullString{String: tmp.Content.Rendered, Valid: true}
	}

	p.Excerpt = sql.NullString{String: "", Valid: false}
	if len(tmp.Excerpt.Rendered) > 0 {
		p.Excerpt = sql.NullString{String: tmp.Excerpt.Rendered, Valid: true}
	}

	// fmt.Println("OG AUTHOR", tmp.Author)

	// Assign other fields
	p.ID = tmp.ID
	p.PostID = tmp.PostID
	p.Date = tmp.Date
	p.Slug = tmp.Slug
	p.Status = tmp.Status
	p.Type = tmp.Type
	p.Link = tmp.Link
	p.Author = tmp.Author
	p.FeaturedMedia = tmp.FeaturedMedia
	// TODO: Assign categories and tags if needed
	// p.Categories = tmp.Categories
	// p.Tags = tmp.Tags
	p.Source = tmp.Source
	return nil
}

// CreatePost inserts a new post into the database.
func CreatePost(db *gorm.DB, post *Post) error {
	result := db.Omit("id").Create(post)

	if result.Error != nil {
		return tracer.TraceError(result.Error)
	}

	return nil
}

// GetPost retrieves a post from the database by ID.
func GetPost(db *gorm.DB, id uint, source string) (Post, error) {
	var post Post
	result := db.Where("post_id=? AND source=?", id, source).First(&post)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return post, errors.New("record not found")
	}

	if result.Error != nil {
		return post, fmt.Errorf("error: %v", result.Error)
	}

	return post, nil
}

// GetPostsByQueryString retrieves all posts that contain queryString in
// Content and Title.
func GetPostsByQueryString(db *gorm.DB, queryString string, page string) ([]Post, error) {
	var posts []Post
	db.Raw(
		`SELECT * 
		FROM posts,
		to_tsvector(posts.title || posts.excerpt) document,
		to_tsquery(@query) query,
		NULLIF(ts_rank(to_tsvector(posts.title), query), 0) rank_title,
    	NULLIF(ts_rank(to_tsvector(posts.excerpt), query), 0) rank_excerpt,
		SIMILARITY(@query, posts.title || posts.excerpt) similarity
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
