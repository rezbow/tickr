package events

import (
	"log/slog"

	"gorm.io/gorm"
)

type EventsService struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewEventsService(db *gorm.DB, logger *slog.Logger) *EventsService {
	return &EventsService{db: db, logger: logger}
}
