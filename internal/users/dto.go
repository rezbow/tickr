package users

import (
	"errors"
	"regexp"

	"github.com/rezbow/tickr/internal/entities"
	"github.com/rezbow/tickr/internal/utils"
)

var ROLES = []string{"admin", "user", "organizer"}

type UserCreateDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func (u *UserCreateDTO) Validate() utils.ValidationErrors {
	validator := utils.NewValidator()
	validator.Must(len(u.Name) > 2 && len(u.Name) < 255, "name", "name must be between 2 and 255 characters")
	validator.Must(len(u.Email) > 2 && len(u.Email) < 255, "email", "email must be between 2 and 255 characters")
	validator.Must(len(u.Password) >= 8, "password", "password must be at least 8 characters")
	validator.In(u.Role, ROLES, "role", "role must be one of admin, user, organizer")
	validator.Regex(u.Email, regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`), "email", "invalid email format")
	if !validator.Valid() {
		return validator.Errors
	}
	return nil
}

type UserUpdateDTO struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Role     *string `json:"role"`
}

func (u *UserUpdateDTO) Validate() utils.ValidationErrors {
	validator := utils.NewValidator()
	if u.Name != nil {
		validator.Must(len(*u.Name) > 2 && len(*u.Name) < 255, "name", "name must be between 2 and 255 characters")
	}
	if u.Email != nil {
		validator.Must(len(*u.Email) > 2 && len(*u.Email) < 255, "email", "email must be between 2 and 255 characters")
		validator.Regex(*u.Email, regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`), "email", "invalid email format")
	}
	if u.Password != nil {
		validator.Must(len(*u.Password) >= 8, "password", "password must be at least 8 characters")
	}
	if u.Role != nil {
		validator.In(*u.Role, ROLES, "role", "role must be one of admin, user, organizer")
	}

	if !validator.Valid() {
		return validator.Errors
	}
	return nil
}

func (u *UserUpdateDTO) ToMap() (map[string]any, error) {
	updates := make(map[string]any)
	if u.Name != nil {
		updates["name"] = *u.Name
	}
	if u.Email != nil {
		updates["email"] = *u.Email
	}
	if u.Password != nil {
		passwordHash, err := hashPassword(*u.Password)
		if err != nil {
			return nil, errors.New("unable to hash the password")
		}
		updates["password_hash"] = passwordHash
	}
	if u.Role != nil {
		updates["role"] = *u.Role
	}

	return updates, nil
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (l *LoginDTO) Validate() utils.ValidationErrors {
	validator := utils.NewValidator()
	validator.Must(len(l.Email) > 0, "email", "email is required")
	validator.Must(len(l.Password) > 0, "password", "password is required")
	validator.Regex(l.Email, regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`), "email", "invalid email format")
	if !validator.Valid() {
		return validator.Errors
	}
	return nil
}

type LoginResponseDTO struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	User         UserResponseDTO `json:"user"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (r *RefreshTokenDTO) Validate() utils.ValidationErrors {
	validator := utils.NewValidator()
	validator.Must(len(r.RefreshToken) > 0, "refresh_token", "refresh token is required")
	if !validator.Valid() {
		return validator.Errors
	}
	return nil
}

type RefreshTokenResponseDTO struct {
	AccessToken string `json:"access_token"`
}

type UserResponseDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func UserEntityToUserResponse(user *entities.User) UserResponseDTO {
	return UserResponseDTO{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

func UserEntitiesToUserResponse(users []entities.User) []UserResponseDTO {
	userResponses := make([]UserResponseDTO, len(users))
	for idx, u := range users {
		userResponses[idx] = UserResponseDTO{
			ID:    u.ID.String(),
			Name:  u.Name,
			Email: u.Email,
			Role:  u.Role,
		}
	}
	return userResponses
}
