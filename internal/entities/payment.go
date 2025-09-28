package entities

import (
	"time"

	"github.com/google/uuid"
)

var (
	PaymentPending   = "pending"
	PaymentConfirmed = "confirmed"
	PaymentCanceled  = "canceled"
)

type Payment struct {
	ID         uuid.UUID
	UserId     uuid.UUID
	TicketId   uuid.UUID
	Quantity   int
	PaidAmount int64
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	// associations
	User   *User   // belongs to
	Ticket *Ticket // belongs to
}
