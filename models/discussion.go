// discussion.go 
package models

import "time"

// Discussion represents a top-level discussion topic.
type Discussion struct {
    ID          int        `json:"id" db:"id"`
    UserID      int        `json:"user_id" db:"user_id"`
    Title       string     `json:"title" db:"title"`
    Content     string     `json:"content" db:"content"`
    ScheduledAt *time.Time `json:"scheduled_at,omitempty" db:"scheduled_at"` // nil ⇒ post immediately
    CreatedAt   time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
