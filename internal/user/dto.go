// dto.go 
package user

import "errors"

// UpdateUserDTO binds JSON for PUT /users/:id.
// All fields are optional; only non‐zero (non‐empty) fields will be updated.
type UpdateUserDTO struct {
    Username *string `json:"username,omitempty"`
    Email    *string `json:"email,omitempty"`
    Password *string `json:"password,omitempty"`
    FullName *string `json:"full_name,omitempty"`
    Bio      *string `json:"bio,omitempty"`
}

// Validate ensures at least one field is present.
func (dto *UpdateUserDTO) Validate() error {
    if dto.Username == nil && dto.Email == nil &&
       dto.Password == nil && dto.FullName == nil && dto.Bio == nil {
        return errors.New("at least one field must be provided")
    }
    return nil
}
