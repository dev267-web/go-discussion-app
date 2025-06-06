// tag.go 
package models

import "time"

// Tag represents a single tag that can be associated with many discussions.
type Tag struct {
    ID        int       `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// DiscussionTag is the join table between discussions and tags.
// If you need to work directly with this table, you can use the struct below.
// Otherwise you can handle many-to-many logic in your repository/service.
type DiscussionTag struct {
    DiscussionID int `json:"discussion_id" db:"discussion_id"`
    TagID        int `json:"tag_id" db:"tag_id"`
}
