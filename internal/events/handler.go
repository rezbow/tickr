package events

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"github.com/rezbow/tickr/internal/utils"
	"gorm.io/gorm"
)

func (service *EventsService) CreateEventHandler(c *gin.Context) {
	var input EventCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors := input.Validate(); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	userIdAny, _ := c.Get("user_id")

	userId, ok := userIdAny.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id "})
	}

	event := &entities.Event{
		Title:     input.Title,
		Venue:     input.Venue,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		UserId:    userId,
	}
	if input.Description != nil {
		event.Description.Valid = true
		event.Description.String = *input.Description
	}

	err := service.createEvent(c.Request.Context(), event)
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		service.logger.Error("failed creating event", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, EventEntityToEventResponse(event))
}

func (service *EventsService) GetEventHandler(c *gin.Context) {
	id := c.Param("id")
	eventId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	event, err := service.getEvent(c.Request.Context(), eventId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		service.logger.Error("failed retrieving event", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, EventEntityToEventResponse(event))

}

func (service *EventsService) GetEventsHandler(c *gin.Context) {
	var p utils.Pagination

	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	events, total, err := service.getEvents(c.Request.Context(), &p)
	if err != nil {
		service.logger.Error("failed to get users", "page", p.Page, "limit", p.PageSize, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":      EventEntitiesToEventResponse(events),
		"total":     total,
		"page":      p.Page,
		"page_size": p.PageSize,
	})
}

func (service *EventsService) DeleteEventHandler(c *gin.Context) {
	id := c.Param("id")
	eventId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	err = service.deleteEvent(c.Request.Context(), eventId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		service.logger.Error("failed retrieving event", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
