package tickets

import (
	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"github.com/rezbow/tickr/internal/utils"
)

type TicketInput struct {
	EventId         uuid.UUID `json:"event_id" binding:"required"`
	Price           int       `json:"price" binding:"required"`
	TotalQuantities int       `json:"total_quantities" binding:"required"`
}

type Ticket struct {
	ID                 uuid.UUID `json:"id"`
	EventId            uuid.UUID `json:"event_id"`
	Price              int       `json:"price"`
	TotalQuantites     int       `json:"total_quantity"`
	RemainingQuantites int       `json:"remaining_quantity"`
}

func TicketEntityToTicket(t *entities.Ticket) Ticket {
	return Ticket{
		ID:                 t.ID,
		EventId:            t.EventId,
		Price:              t.Price,
		TotalQuantites:     t.TotalQuantities,
		RemainingQuantites: t.RemainingQuantities,
	}
}

func (t *TicketInput) Validate() utils.ValidationErrors {
	v := utils.NewValidator()
	v.Must(t.Price > 0, "price", "must be positive integer")
	v.Must(t.TotalQuantities > 0, "total_quantities", "must be positive integer")
	if !v.Valid() {
		return v.Errors
	}
	return nil
}
