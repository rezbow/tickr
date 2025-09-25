package users

import (
	"log/slog"

	"github.com/rezbow/tickr/internal/auth"
	"gorm.io/gorm"
)

type UsersService struct {
	db                 *gorm.DB
	logger             *slog.Logger
	jwtService         *auth.JWTService
	refreshTokenService *auth.RefreshTokenService
}

func NewUserService(db *gorm.DB, logger *slog.Logger) *UsersService {
	return &UsersService{
		db:                 db,
		logger:             logger,
		jwtService:         auth.NewJWTService(),
		refreshTokenService: auth.NewRefreshTokenService(db, logger),
	}
}
