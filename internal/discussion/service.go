// service.go 
package discussion

import (
    "context"
    "time"

    "go-discussion-app/models"
		tagpkg "go-discussion-app/internal/tag"
)

type Service interface {
    Create(ctx context.Context, userID int, dto *CreateDiscussionDTO) (int, error)
    GetAll(ctx context.Context) ([]models.Discussion, error)
    GetByID(ctx context.Context, id int) (*models.Discussion, error)
    Update(ctx context.Context, id int, dto *UpdateDiscussionDTO) (*models.Discussion, error)
    Delete(ctx context.Context, id int) error

    GetByUser(ctx context.Context, userID int) ([]models.Discussion, error)
    GetByTag(ctx context.Context, tag string) ([]models.Discussion, error)
    AddTags(ctx context.Context, discussionID int, dto *AddTagsDTO) error
    Schedule(ctx context.Context, userID int, dto *ScheduleDTO) (int, error)
}

type service struct {
    repo    Repository
    tagRepo tagpkg.TagRepository
}

func NewService(
    repo Repository,
    tagRepo tagpkg.TagRepository,
) Service {
    return &service{repo: repo, tagRepo: tagRepo}
}


func (s *service) Create(ctx context.Context, userID int, dto *CreateDiscussionDTO) (int, error) {
    d := &models.Discussion{
        UserID:      userID,
        Title:       dto.Title,
        Content:     dto.Content,
        ScheduledAt: dto.ScheduledAt,
        CreatedAt:   time.Now().UTC(),
        UpdatedAt:   time.Now().UTC(),
    }
    return s.repo.Create(ctx, d)
}

func (s *service) GetAll(ctx context.Context) ([]models.Discussion, error) {
    return s.repo.GetAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id int) (*models.Discussion, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id int, dto *UpdateDiscussionDTO) (*models.Discussion, error) {
    d, err := s.repo.GetByID(ctx, id)
    if err != nil || d == nil {
        return nil, err
    }
    if dto.Title != nil {
        d.Title = *dto.Title
    }
    if dto.Content != nil {
        d.Content = *dto.Content
    }
    if dto.ScheduledAt != nil {
        d.ScheduledAt = dto.ScheduledAt
    }
    d.UpdatedAt = time.Now().UTC()
    if err := s.repo.Update(ctx, d); err != nil {
        return nil, err
    }
    return d, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
    return s.repo.Delete(ctx, id)
}

func (s *service) GetByUser(ctx context.Context, userID int) ([]models.Discussion, error) {
    return s.repo.GetByUser(ctx, userID)
}

func (s *service) GetByTag(ctx context.Context, tag string) ([]models.Discussion, error) {
    return s.repo.GetByTag(ctx, tag)
}

func (s *service) AddTags(
    ctx context.Context,
    discussionID int,
    dto *AddTagsDTO,
) error {
    // Gather tag IDs, creating tags if they do not exist
    var tagIDs []int
    for _, name := range dto.Tags {
        t, err := s.tagRepo.GetByName(ctx, name)
        if err != nil {
            return err
        }
        if t == nil {
            // Tag doesn’t exist → create it
            newID, err := s.tagRepo.Create(ctx, name)
            if err != nil {
                return err
            }
            tagIDs = append(tagIDs, newID)
        } else {
            tagIDs = append(tagIDs, t.ID)
        }
    }

    // Delegate to discussion_tags join table insertion
    return s.repo.AddTags(ctx, discussionID, tagIDs)
}

func (s *service) Schedule(ctx context.Context, userID int, dto *ScheduleDTO) (int, error) {
    d := &models.Discussion{
        UserID:      userID,
        Title:       dto.Title,
        Content:     dto.Content,
        ScheduledAt: &dto.ScheduledAt,
        CreatedAt:   time.Now().UTC(),
        UpdatedAt:   time.Now().UTC(),
    }
    return s.repo.Create(ctx, d)
}
