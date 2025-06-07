package subscription

import "time"

type SubscribeDTO struct {
	Email        string    `json:"email" binding:"required,email"`
	SubscribedAt time.Time `json:"subscribed_at" binding:"required"`
}
