// repository.go 
package subscription

import (
	"database/sql"
	"go-discussion-app/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateSubscription(sub *models.Subscription) error {
	query := `INSERT INTO subscriptions (discussion_id, user_id, email, subscribed_at)
	          VALUES ($1, $2, $3, $4)
			  ON CONFLICT (discussion_id, email) DO NOTHING`
	_, err := r.db.Exec(query, sub.DiscussionID, sub.UserID, sub.Email, sub.SubscribedAt)
	return err
}

func (r *Repository) DeleteSubscription(discussionID int, email string) error {
	query := `DELETE FROM subscriptions WHERE discussion_id = $1 AND email = $2`
	_, err := r.db.Exec(query, discussionID, email)
	return err
}

func (r *Repository) GetSubscriberEmails(discussionID int) ([]string, error) {
	rows, err := r.db.Query(`SELECT email FROM subscriptions WHERE discussion_id = $1`, discussionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}
