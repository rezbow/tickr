package payment

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (service *PaymentService) GetPaymentHandler(c *gin.Context) {
	idStr := c.Param("id")
	paymentId, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	payment, err := service.getPayment(c.Request.Context(), paymentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		service.logger.Error("Failed getting payment", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, PaymentEntityToPayment(*payment))
}

func (service *PaymentService) BuyTicketHandler(c *gin.Context) {
	ErrInsuffcientQuantity := errors.New("insufficient quantities")
	var paymentDetail PaymentDetail
	if err := c.ShouldBindJSON(&paymentDetail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors := paymentDetail.Validate(); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}
	var paymentId uuid.UUID
	err := service.db.Transaction(func(tx *gorm.DB) error {
		var ticket entities.Ticket
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", paymentDetail.TicketId).First(&ticket).Error; err != nil {
			return err
		}

		if ticket.RemainingQuantities == 0 || ticket.RemainingQuantities < paymentDetail.Quantity {
			return ErrInsuffcientQuantity
		}

		ticket.RemainingQuantities -= paymentDetail.Quantity
		if err := service.db.Save(&ticket).Error; err != nil {
			return err
		}

		payment := entities.Payment{
			ID:         uuid.New(),
			UserId:     paymentDetail.UserId,
			TicketId:   paymentDetail.TicketId,
			Quantity:   paymentDetail.Quantity,
			PaidAmount: paymentDetail.Quantity * ticket.Price,
		}
		if err := service.db.Create(&payment).Error; err != nil {
			return err
		}
		paymentId = payment.ID
		return nil
	})
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		case gorm.ErrForeignKeyViolated:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user or ticket "})
		case ErrInsuffcientQuantity:
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient quantity"})
		default:
			service.logger.Error("payment failed", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"payment_id": paymentId, "message": "payment successful"})
}
