// dto.go 
package auth

import "errors"

// RegisterDTO is the payload for POST /auth/register
type RegisterDTO struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    FullName string `json:"full_name,omitempty"`
    Bio      string `json:"bio,omitempty"`
}

func (dto *RegisterDTO) Validate() error {
    if dto.Username == "" {
        return errors.New("username is required")
    }
    if dto.Email == "" {
        return errors.New("email is required")
    }
    if dto.Password == "" {
        return errors.New("password is required")
    }
    return nil
}

// LoginDTO is the payload for POST /auth/login
type LoginDTO struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (dto *LoginDTO) Validate() error {
    if dto.Email == "" {
        return errors.New("email is required")
    }
    if dto.Password == "" {
        return errors.New("password is required")
    }
    return nil
}
