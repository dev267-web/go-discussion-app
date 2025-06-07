// service.go 
package comment

import (
    "context"
    "time"

    "go-discussion-app/models"
)

type Service interface {
    AddComment(ctx context.Context, discussionID, userID int, content string) (int, error)
    GetComments(ctx context.Context, discussionID int) ([]models.Comment, error)
}

type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
}

func (s *service) AddComment(ctx context.Context, discussionID, userID int, content string) (int, error) {
    comment := &models.Comment{
        DiscussionID: discussionID,
        UserID:       userID,
        Content:      content,
        CreatedAt:    time.Now().UTC(),
    }
    return s.repo.Create(ctx, comment)
}

func (s *service) GetComments(ctx context.Context, discussionID int) ([]models.Comment, error) {
    return s.repo.ListByDiscussion(ctx, discussionID)
}
