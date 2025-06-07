// service.go 
package subscription

import (
	"fmt"
	"go-discussion-app/models"
	"go-discussion-app/pkg/mailer"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) Subscribe(sub *models.Subscription) error {
	return s.repo.CreateSubscription(sub)
}

func (s *Service) Unsubscribe(discussionID int, email string) error {
	return s.repo.DeleteSubscription(discussionID, email)
}

func (s *Service) NotifySubscribers(discussionID int, subject, body string) error {
	emails, err := s.repo.GetSubscriberEmails(discussionID)
	if err != nil {
		return fmt.Errorf("failed to get emails: %w", err)
	}
	if len(emails) == 0 {
		return nil
	}
	return mailer.SendMail(emails, subject, body)
}
