// service.go 
package auth

import (
    "context"
    "errors"
    "time"

    "golang.org/x/crypto/bcrypt"

    "go-discussion-app/internal/user"
    "go-discussion-app/models"
    "go-discussion-app/pkg/jwtutil"
)

var (
    ErrUserExists         = errors.New("user with that email already exists")
    ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService struct {
    userRepo user.UserRepository
}

func NewService(uRepo user.UserRepository) *AuthService {
    return &AuthService{userRepo: uRepo}
}

func (s *AuthService) Register(ctx context.Context, dto *RegisterDTO) (int, error) {
    if err := dto.Validate(); err != nil {
        return 0, err
    }

    if existing, err := s.userRepo.GetByEmail(ctx, dto.Email); err != nil {
        return 0, err
    } else if existing != nil {
        return 0, ErrUserExists
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
    if err != nil {
        return 0, err
    }

    now := time.Now().UTC()
    u := &models.User{
        Username:     dto.Username,
        Email:        dto.Email,
        PasswordHash: string(hashed),
        FullName:     dto.FullName,
        Bio:          dto.Bio,
        CreatedAt:    now,
        UpdatedAt:    now,
    }
    return s.userRepo.Create(ctx, u)
}

func (s *AuthService) Login(ctx context.Context, dto *LoginDTO) (string, error) {
    if err := dto.Validate(); err != nil {
        return "", err
    }

    u, err := s.userRepo.GetByEmail(ctx, dto.Email)
    if err != nil {
        return "", err
    }
    if u == nil {
        return "", ErrInvalidCredentials
    }
    if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(dto.Password)); err != nil {
        return "", ErrInvalidCredentials
    }

    return jwtutil.GenerateToken(u.ID)
}
