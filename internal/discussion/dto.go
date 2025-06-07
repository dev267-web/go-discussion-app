// dto.go 
package discussion

import (
    "errors"
    "time"
)

// CreateDiscussionDTO for POST /discussions
type CreateDiscussionDTO struct {
    Title       string     `json:"title"`
    Content     string     `json:"content"`
    ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
}

func (dto *CreateDiscussionDTO) Validate() error {
    if dto.Title == "" {
        return errors.New("title is required")
    }
    if dto.Content == "" {
        return errors.New("content is required")
    }
    return nil
}

// UpdateDiscussionDTO for PUT /discussions/:id
type UpdateDiscussionDTO struct {
    Title       *string    `json:"title,omitempty"`
    Content     *string    `json:"content,omitempty"`
    ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
}

func (dto *UpdateDiscussionDTO) Validate() error {
    if dto.Title == nil && dto.Content == nil && dto.ScheduledAt == nil {
        return errors.New("at least one field must be provided")
    }
    return nil
}

// AddTagsDTO for POST /discussions/:id/tags
type AddTagsDTO struct {
    Tags []string `json:"tags"` // tag names
}

func (dto *AddTagsDTO) Validate() error {
    if len(dto.Tags) == 0 {
        return errors.New("tags list cannot be empty")
    }
    return nil
}

// ScheduleDTO for POST /discussions/schedule
type ScheduleDTO struct {
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    ScheduledAt time.Time `json:"scheduled_at"`
}

func (dto *ScheduleDTO) Validate() error {
    if dto.Title == "" {
        return errors.New("title is required")
    }
    if dto.Content == "" {
        return errors.New("content is required")
    }
    if dto.ScheduledAt.IsZero() {
        return errors.New("scheduled_at is required")
    }
    return nil
}
