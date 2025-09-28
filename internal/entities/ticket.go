package entities

import (
	"time"

	"github.com/google/uuid"
)

// gorm model
type Ticket struct {
	ID                  uuid.UUID
	EventId             uuid.UUID
	UserId              uuid.UUID
	Price               int64
	TotalQuantities     int
	RemainingQuantities int
	CreatedAt           time.Time
	UpdatedAt           time.Time
	// associations
	Event Event
}
