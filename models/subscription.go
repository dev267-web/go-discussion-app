// subscription.go 
package models

import "time"

// Subscription represents an email subscription for a discussion.
type Subscription struct {
    ID           int       `json:"id" db:"id"`
    DiscussionID int       `json:"discussion_id" db:"discussion_id"`
    UserID       *int      `json:"user_id,omitempty" db:"user_id"` // nullable; stored as NULL if external email
    Email        string    `json:"email" db:"email"`
    SubscribedAt time.Time `json:"subscribed_at" db:"subscribed_at"`
}
