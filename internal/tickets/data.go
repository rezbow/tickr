package tickets

import (
	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"github.com/rezbow/tickr/internal/utils"
)

type TicketInput struct {
	EventId         uuid.UUID `json:"event_id" binding:"required"`
	Price           int64       `json:"price" binding:"required"`
	TotalQuantities int       `json:"total_quantities" binding:"required"`
}

type Ticket struct {
	ID                 uuid.UUID `json:"id"`
	EventId            uuid.UUID `json:"event_id"`
	Price              int64       `json:"price"`
    TotalQuantities    int       `json:"total_quantities"`
    RemainingQuantities int      `json:"remaining_quantities"`
}

func TicketEntityToTicket(t *entities.Ticket) Ticket {
	return Ticket{
		ID:                 t.ID,
		EventId:            t.EventId,
		Price:              t.Price,
        TotalQuantities:    t.TotalQuantities,
        RemainingQuantities: t.RemainingQuantities,
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
