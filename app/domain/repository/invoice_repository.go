package repository

import (
	"github.com/ijufumi/practice-202512/app/domain/models"
	"time"

	"gorm.io/gorm"
)

type InvoiceRepository interface {
	Create(db *gorm.DB, invoice *models.Invoice) error
	FindByPaymentDueDateRange(db *gorm.DB, startDate, endDate *time.Time, offset, limit int) ([]*models.Invoice, error)
}
