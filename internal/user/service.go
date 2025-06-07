// service.go 
package user

import (
    "context"
    //"database/sql"
    "errors"
    "time"

    "golang.org/x/crypto/bcrypt"
    "go-discussion-app/models"
)

var (
    ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
    repo UserRepository
}

func NewService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

// GetByID fetches a user by ID.
func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
    u, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if u == nil {
        return nil, ErrUserNotFound
    }
    return u, nil
}

// Update applies non‐nil fields from dto to the existing user.
func (s *UserService) Update(ctx context.Context, id int, dto *UpdateUserDTO) (*models.User, error) {
    existing, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if existing == nil {
        return nil, ErrUserNotFound
    }

    // Apply updates
    if dto.Username != nil {
        existing.Username = *dto.Username
    }
    if dto.Email != nil {
        existing.Email = *dto.Email
    }
    if dto.Password != nil {
        hashed, err := bcrypt.GenerateFromPassword([]byte(*dto.Password), bcrypt.DefaultCost)
        if err != nil {
            return nil, err
        }
        existing.PasswordHash = string(hashed)
    }
    if dto.FullName != nil {
        existing.FullName = *dto.FullName
    }
    if dto.Bio != nil {
        existing.Bio = *dto.Bio
    }
    existing.UpdatedAt = time.Now().UTC()

    if _, err := s.repo.Update(ctx, existing); err != nil {
        return nil, err
    }
    return existing, nil
}

// Delete removes a user by ID.
func (s *UserService) Delete(ctx context.Context, id int) error {
    // Optionally, check existence first:
    u, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if u == nil {
        return ErrUserNotFound
    }
    _, err = s.repo.Delete(ctx, id)
    return err
}
