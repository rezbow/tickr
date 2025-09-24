package events

import (
	"context"

	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"gorm.io/gorm"
)

func (service *EventsService) createEvent(ctx context.Context, event *entities.Event) error {

	event.ID = uuid.New()
	err := gorm.G[entities.Event](service.db).Create(ctx, event)
	if err != nil {
		return err
	}
	return nil
}

func (service *EventsService) getEvent(ctx context.Context, eventId uuid.UUID) (*entities.Event, error) {
	event, err := gorm.G[entities.Event](service.db).Where("id = ?", eventId).First(ctx)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (service *EventsService) deleteEvent(ctx context.Context, eventId uuid.UUID) error {
	rowsAffected, err := gorm.G[entities.Event](service.db).Where("id = ?", eventId).Delete(ctx)
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	} else if err != nil {
		return err
	}
	return nil
}

func (service *EventsService) getEvents(ctx context.Context, page, limit int) ([]entities.Event, int64, error) {
	var total int64
	if res := service.db.Model(&entities.Event{}).Count(&total); res.Error != nil {
		return nil, 0, res.Error
	}
	users, err := gorm.G[entities.Event](service.db).Offset((page - 1) * limit).Limit(limit).Find(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
