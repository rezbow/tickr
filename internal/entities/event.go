package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// gorm model
type Event struct {
	ID          uuid.UUID
	Title       string
	Description sql.NullString
	Venue       string
	UserId      uuid.UUID
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	// associations
	User    User     // Belongs to
	Tickets []Ticket // has many
}
