package payment

import (
	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"github.com/rezbow/tickr/internal/utils"
)

type PaymentDetail struct {
	TicketId uuid.UUID `json:"ticket_id" binding:"required"`
	UserId   uuid.UUID `json:"user_id" binding:"required"`
	Quantity int       `json:"quantity" binding:"required"`
}

type Payment struct {
	ID       uuid.UUID `json:"id"`
	UserId   uuid.UUID `json:"user_id"`
	Ticket   uuid.UUID `json:"ticket_id"`
	Quantity int       `json:"quantity"`
}

func PaymentEntityToPayment(p entities.Payment) Payment {
	return Payment{
		ID:       p.ID,
		UserId:   p.UserId,
		Ticket:   p.TicketId,
		Quantity: p.Quantity,
	}
}

func (pd *PaymentDetail) Validate() utils.ValidationErrors {
	validator := utils.NewValidator()
	validator.Must(pd.Quantity > 0, "quantity", "Quantity must be greater than 0")
	return validator.Errors
}
