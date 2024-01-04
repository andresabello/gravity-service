package posts

import (
	"fmt"
	"pi-search/pkg/tracer"

	"gorm.io/gorm"
)

func ingestUsingQueue(db *gorm.DB, post *Post) error {
	foundPost, err := GetPost(db, post.PostID, post.Source)
	if err != nil {
		// Print error. Not necesary to take it as a real error as this
		// indicates no posts were found so it will skip to update.
		fmt.Println(err)
		err := CreatePost(db, post)
		if err != nil {
			return tracer.TraceError(
				fmt.Errorf(
					"unable to create post from database %d error %s",
					post.PostID,
					err,
				),
			)
		}

		return nil
	}

	err = UpdatePost(db, &foundPost)
	if err != nil {
		return tracer.TraceError(
			fmt.Errorf(
				"unable to update post from database %d error %s",
				post.ID,
				err,
			),
		)
	}

	return nil
}
