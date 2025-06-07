package comment

import "errors"

// CreateCommentDTO binds the JSON body for creating a comment.
type CreateCommentDTO struct {
    Content string `json:"content"`
}

// Validate ensures the content is not empty.
func (dto *CreateCommentDTO) Validate() error {
    if dto.Content == "" {
        return errors.New("content is required")
    }
    return nil
}
