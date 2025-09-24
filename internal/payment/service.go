package payment

import (
	"log/slog"

	"gorm.io/gorm"
)

type PaymentService struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewPaymentService(db *gorm.DB, logger *slog.Logger) *PaymentService {
	return &PaymentService{db: db, logger: logger}
}
