package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"github.com/rezbow/tickr/internal/utils"
)

func EventEntitiesToEventResponse(events []entities.Event) []EventResponseDTO {
	result := make([]EventResponseDTO, len(events))
	for i, e := range events {
		result[i] = EventEntityToEventResponse(&e)
	}
	return result
}

type EventCreateDTO struct {
	Title       string    `json:"title" binding:"required"`
	Description *string   `json:"description"`
	Venue       string    `json:"venue" binding:"required"`
	UserId      uuid.UUID `json:"user_id" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
}

func (e *EventCreateDTO) Validate() utils.ValidationErrors {
	validator := utils.NewValidator()

	validator.Must(len(e.Title) >= 2 && len(e.Title) <= 255, "title", "title must be between 2 and 255 characters")
	if e.Description != nil {
		validator.Must(len(*e.Description) >= 2 && len(*e.Description) <= 1024, "description", "title must be between 2 and 1024 characters")
	}
	validator.Must(len(e.Venue) >= 2 && len(e.Venue) <= 255, "title", "title must be between 2 and 255 characters")
	// start time, end time
	validator.Must(e.StartTime.After(time.Now()), "start_time", "start_time should be in future")
	validator.Must(e.EndTime.After(time.Now()), "end_time", "end_time should be in future")
	validator.Must(e.EndTime.After(e.StartTime), "end_time", "end_time should be after start_time ")

	if !validator.Valid() {
		return validator.Errors
	}
	return nil
}

type EventUpdateDTO struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Venue       *string    `json:"venue"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
}

type EventResponseDTO struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Venue       string    `json:"venue"`
	UserId      uuid.UUID `json:"user_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func EventEntityToEventResponse(e *entities.Event) EventResponseDTO {
	return EventResponseDTO{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description.String,
		Venue:       e.Venue,
		UserId:      e.UserId,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// -------------------------------------------------------- //

type TicketInput struct {
	EventId         uuid.UUID `json:"event_id" binding:"required"`
	Price           int       `json:"price" binding:"required"`
	TotalQuantities int       `json:"total_quantities" binding:"required"`
}

func (t *TicketInput) Validate() utils.ValidationErrors {
	v := utils.NewValidator()
	v.Must(t.Price > 0, "price", "must be positive integer")
	v.Must(t.TotalQuantities > 0, "total_quantities", "must be positive integer")
	return v.Errors
}
