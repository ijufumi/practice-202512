package gateway

import (
	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/domain/repository"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/dao"
	"time"

	"gorm.io/gorm"
)

type invoiceRepository struct{}

func NewInvoiceRepository() repository.InvoiceRepository {
	return &invoiceRepository{}
}

func (r *invoiceRepository) Create(db *gorm.DB, invoice *models.Invoice) error {
	daoInvoice := invoice.ToDAO()
	if err := db.Create(&daoInvoice).Error; err != nil {
		return err
	}
	invoice.ID = daoInvoice.ID
	invoice.CreatedAt = daoInvoice.CreatedAt
	invoice.UpdatedAt = daoInvoice.UpdatedAt
	return nil
}

func (r *invoiceRepository) FindByPaymentDueDateRange(db *gorm.DB, startDate, endDate *time.Time, offset, limit int) ([]*models.Invoice, error) {
	var daoInvoices []*dao.Invoice
	var daoStartDate, daoEndDate time.Time
	if startDate != nil {
		daoStartDate = *startDate
	} else {
		daoStartDate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
	}
	if endDate != nil {
		daoEndDate = *endDate
	} else {
		daoEndDate = time.Date(9999, 12, 31, 23, 59, 59, 999999, time.Local)
	}
	if err := db.Where("payment_due_date BETWEEN ? AND ?", daoStartDate, daoEndDate).
		Order("payment_due_date ASC").
		Offset(offset).
		Limit(limit).
		Find(&daoInvoices).Error; err != nil {
		return nil, err
	}

	invoices := make([]*models.Invoice, len(daoInvoices))
	for i, daoInvoice := range daoInvoices {
		invoices[i] = models.InvoiceFromDAO(daoInvoice)
	}
	return invoices, nil
}
