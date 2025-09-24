package tickets

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

func (service *TicketsService) CreateTicket(c *gin.Context) {
	var input TicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors := input.Validate(); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	ticket := entities.Ticket{
		EventId:         input.EventId,
		Price:           input.Price,
		TotalQuantities: input.TotalQuantities,
	}

	if err := service.createTicket(c.Request.Context(), &ticket); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
			return
		}
		service.logger.Error("failed to create ticket", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "ticket created"})
}

func (service *TicketsService) GetTicket(c *gin.Context) {
	id := c.Param("id")
	ticketId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}
	ticket, err := service.getTicket(c.Request.Context(), ticketId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}
		service.logger.Error("failed to get ticket", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, TicketEntityToTicket(ticket))
}

func (service *TicketsService) DeleteTicket(c *gin.Context) {
	id := c.Param("id")
	ticketId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}
	if err := service.deleteTicket(c.Request.Context(), ticketId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}
		service.logger.Error("failed to delete ticket", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ticket deleted"})
}

func (service *TicketsService) GetEventTicketsHandler(c *gin.Context) {
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

/*
func (service *TicketsService) CreateTicketForEvent(c *gin.Context) {
	idStr := c.Param("id")
	eventId, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
	}

	var input TicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors := input.Validate(); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors})
		return
	}

	ticket := entities.Ticket{
		EventId:         input.EventId,
		Price:           input.Price,
		TotalQuantities: input.TotalQuantities,
	}

	if err := service.createTicket(c.Request.Context(), &ticket); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
			return
		}
		service.logger.Error("failed to create ticket", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "ticket created"})
}
*/
