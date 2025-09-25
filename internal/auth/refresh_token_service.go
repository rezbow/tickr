package auth

import (
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"gorm.io/gorm"
)

type RefreshTokenService struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewRefreshTokenService(db *gorm.DB, logger *slog.Logger) *RefreshTokenService {
	return &RefreshTokenService{
		db:     db,
		logger: logger,
	}
}

func (r *RefreshTokenService) CreateRefreshToken(userID uuid.UUID, token string) (*entities.RefreshToken, error) {
	refreshToken := &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := r.db.Create(refreshToken).Error; err != nil {
		r.logger.Error("failed to create refresh token", "userId", userID.String(), "error", err)
		return nil, err
	}

	return refreshToken, nil
}

func (r *RefreshTokenService) GetRefreshToken(token string) (*entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken
	if err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found or expired")
		}
		r.logger.Error("failed to get refresh token", "error", err)
		return nil, err
	}

	return &refreshToken, nil
}

func (r *RefreshTokenService) DeleteRefreshToken(token string) error {
	if err := r.db.Where("token = ?", token).Delete(&entities.RefreshToken{}).Error; err != nil {
		r.logger.Error("failed to delete refresh token", "error", err)
		return err
	}
	return nil
}

func (r *RefreshTokenService) DeleteAllUserRefreshTokens(userID uuid.UUID) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&entities.RefreshToken{}).Error; err != nil {
		r.logger.Error("failed to delete user refresh tokens", "userId", userID.String(), "error", err)
		return err
	}
	return nil
}

func (r *RefreshTokenService) CleanupExpiredTokens() error {
	if err := r.db.Where("expires_at < ?", time.Now()).Delete(&entities.RefreshToken{}).Error; err != nil {
		r.logger.Error("failed to cleanup expired refresh tokens", "error", err)
		return err
	}
	return nil
}