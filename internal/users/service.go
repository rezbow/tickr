package users

import (
	"log/slog"

	"gorm.io/gorm"
)

type UsersService struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserService(db *gorm.DB, logger *slog.Logger) *UsersService {
	return &UsersService{db: db, logger: logger}
}
