package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Price int64

func (p Price) toString() string {
	dollars := p / 100
	cents := p % 100
	return fmt.Sprintf("%d.%02d", dollars, cents)
}

// gorm model
type Ticket struct {
	ID                uuid.UUID
	EventId           uuid.UUID
	Price             Price
	TotalQuantity     int
	RemainingQuantity int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	// associations
	// Event belongs to User
	Event Event
}
