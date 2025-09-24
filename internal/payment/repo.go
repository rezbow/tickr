package payment

import (
	"context"

	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"gorm.io/gorm"
)

func (service *PaymentService) createPayment(ctx context.Context, payment *entities.Payment) error {
	err := gorm.G[entities.Payment](service.db).Create(ctx, payment)
	if err != nil {
		return err
	}
	return nil
}

func (service *PaymentService) getPayment(ctx context.Context, paymentId uuid.UUID) (*entities.Payment, error) {
	payment, err := gorm.G[entities.Payment](service.db).Where("id = ?", paymentId).First(ctx)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (service *PaymentService) getTicket(ctx context.Context, ticketId uuid.UUID) (*entities.Ticket, error) {
	ticket, err := gorm.G[entities.Ticket](service.db).Where("id = ?", ticketId).First(ctx)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}
