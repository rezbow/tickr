package tickets

import (
	"log/slog"

	"gorm.io/gorm"
)

type TicketsService struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewTicketsService(db *gorm.DB, logger *slog.Logger) *TicketsService {
	return &TicketsService{db: db, logger: logger}
}
