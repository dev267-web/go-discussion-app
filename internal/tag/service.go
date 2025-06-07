// service.go 
package tag

import (
    "context"

    "go-discussion-app/models"
)

// TagService provides tag‐related business logic.
type TagService struct {
    repo TagRepository
}

// NewService constructs a TagService.
func NewService(repo TagRepository) *TagService {
    return &TagService{repo: repo}
}

// ListTags returns all available tags.
func (s *TagService) ListTags(ctx context.Context) ([]models.Tag, error) {
    return s.repo.GetAll(ctx)
}
