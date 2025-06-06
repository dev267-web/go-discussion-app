// comment.go 
package models

import "time"

// Comment represents a user’s comment on a discussion.
type Comment struct {
    ID           int       `json:"id" db:"id"`
    DiscussionID int       `json:"discussion_id" db:"discussion_id"`
    UserID       int       `json:"user_id" db:"user_id"`
    Content      string    `json:"content" db:"content"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
