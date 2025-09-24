package tickets

import (
	"context"

	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"gorm.io/gorm"
)

func (service *TicketsService) createTicket(ctx context.Context, ticket *entities.Ticket) error {
	ticket.ID = uuid.New()
	err := gorm.G[entities.Ticket](service.db).Create(ctx, ticket)
	if err != nil {
		return err
	}
	return nil
}

func (service *TicketsService) getTicket(ctx context.Context, id uuid.UUID) (*entities.Ticket, error) {
	var ticket entities.Ticket
	ticket, err := gorm.G[entities.Ticket](service.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (service *TicketsService) deleteTicket(ctx context.Context, id uuid.UUID) error {
	rowsAffected, err := gorm.G[entities.Ticket](service.db).Where("id = ?", id).Delete(ctx)
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	} else if err != nil {
		return err
	}
	return nil
}

func (service *TicketsService) getEventTickets(ctx context.Context, eventId uuid.UUID, page, limit int) ([]entities.Ticket, int64, error) {
	var total int64
	if res := service.db.Model(&entities.Ticket{}).Where("event_id", eventId).Count(&total); res.Error != nil {
		return nil, 0, res.Error
	}
	tickets, err := gorm.G[entities.Ticket](service.db).Where("event_id", eventId).Offset((page - 1) * limit).Limit(limit).Find(ctx)
	if err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}
