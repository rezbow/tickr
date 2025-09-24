package tickets

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"github.com/rezbow/tickr/internal/utils"
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
	var p utils.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	id := c.Param("id")
	eventId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	tickets, total, err := service.getEventTickets(c.Request.Context(), eventId, &p)
	if err != nil {
		service.logger.Error("failed to get users", "page", p.Page, "limit", p.PageSize, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      tickets,
		"total":     total,
		"page":      p.Page,
		"page_size": p.PageSize,
	})
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
