package payment

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrInsuffcientQuantity = errors.New("insufficient quantities")

type PaymentService struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewPaymentService(db *gorm.DB, logger *slog.Logger) *PaymentService {
	return &PaymentService{db: db, logger: logger}
}

func (svc *PaymentService) createPayment(p PaymentDetail) (*entities.Payment, error) {
	var payment entities.Payment
	err := svc.db.Transaction(func(tx *gorm.DB) error {
		var ticket entities.Ticket
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", p.TicketId).First(&ticket).Error; err != nil {
			return err
		}

		if ticket.RemainingQuantities == 0 || ticket.RemainingQuantities < p.Quantity {
			return ErrInsuffcientQuantity
		}

		ticket.RemainingQuantities -= p.Quantity
		if err := tx.Save(&ticket).Error; err != nil {
			return err
		}

		payment := entities.Payment{
			ID:         uuid.New(),
			UserId:     p.UserId,
			TicketId:   p.TicketId,
			Quantity:   p.Quantity,
			Status:     entities.PaymentConfirmed,
			PaidAmount: int64(p.Quantity) * ticket.Price,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", payment.ID).First(&payment).Error; err != nil {
			return err
		}
		return nil
	})
	return &payment, err
}
