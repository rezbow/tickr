package payment

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	var paymentDetail PaymentDetail
	if err := c.ShouldBindJSON(&paymentDetail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors := paymentDetail.Validate(); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}
	payment, err := service.createPayment(paymentDetail)
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

	c.JSON(http.StatusCreated, PaymentEntityToPayment(*payment))
}
