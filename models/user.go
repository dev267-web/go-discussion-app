// user.go 
package models

import "time"

// User represents a registered user / profile.
type User struct {
    ID           int       `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"-" db:"password_hash"` // omit hash from JSON responses
    FullName     string    `json:"full_name,omitempty" db:"full_name"`
    Bio          string    `json:"bio,omitempty" db:"bio"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
