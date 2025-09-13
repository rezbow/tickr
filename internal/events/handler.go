package events

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"gorm.io/gorm"
)

func (service *EventsService) CreateEventHandler(c *gin.Context) {
	var input EventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors := input.Validate(); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	event := &entities.Event{
		Title:     input.Title,
		Venue:     input.Venue,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		UserId:    input.UserId,
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

	c.JSON(http.StatusOK, EventToEventRepr(event))
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

	c.JSON(http.StatusOK, EventToEventRepr(event))

}

func (service *EventsService) GetEventsHandler(c *gin.Context) {
	var (
		page  int
		limit int
		err   error
	)

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err = strconv.Atoi(pageStr)
	if page <= 0 || err != nil {
		page = 1
	}
	limit, err = strconv.Atoi(limitStr)
	if limit <= 0 || err != nil {
		limit = 10
	}

	events, total, err := service.getEvents(c.Request.Context(), page, limit)
	if err != nil {
		service.logger.Error("failed to get users", "page", page, "limit", limit, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response struct {
		Data      []EventRepr `json:"data"`
		Page      int         `json:"page"`
		Limit     int         `json:"limit"`
		Total     int64       `json:"total"`
		TotalPage int         `json:"total_page"`
	}
	response.Data = EventsToRepr(events)
	response.Page = page
	response.Limit = limit
	response.Total = total
	response.TotalPage = int(math.Ceil(float64(total) / float64(limit)))
	c.JSON(http.StatusOK, response)

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

	c.JSON(http.StatusNoContent, nil)
}

func (service *EventsService) GetEventTicketsHandler(c *gin.Context) {
	var (
		page  int
		limit int
		err   error
	)

	id := c.Param("id")
	eventId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err = strconv.Atoi(pageStr)
	if page <= 0 || err != nil {
		page = 1
	}
	limit, err = strconv.Atoi(limitStr)
	if limit <= 0 || err != nil {
		limit = 10
	}

	tickets, total, err := service.getEventTickets(c.Request.Context(), eventId, page, limit)
	if err != nil {
		service.logger.Error("failed to get users", "page", page, "limit", limit, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response struct {
		Data      []entities.Ticket `json:"data"`
		Page      int               `json:"page"`
		Limit     int               `json:"limit"`
		Total     int64             `json:"total"`
		TotalPage int               `json:"total_page"`
	}

	response.Data = tickets
	response.Page = page
	response.Limit = limit
	response.Total = total
	response.TotalPage = int(math.Ceil(float64(total) / float64(limit)))
	c.JSON(http.StatusOK, response)

}

func (service *EventsService) UpdateEventHandler(c *gin.Context) {}
