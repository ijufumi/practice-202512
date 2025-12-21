package models

import (
	"github.com/ijufumi/practice-202512/app/domain/value"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/dao"

	"time"
)

type Invoice struct {
	ID             string
	CompanyID      string
	ClientID       string
	IssueDate      time.Time
	PaymentAmount  int
	Fee            int
	FeeRate        float64
	Tax            int
	TaxRate        float64
	InvoiceAmount  int
	PaymentDueDate time.Time
	Status         value.InvoiceStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (i *Invoice) ToDAO() *dao.Invoice {
	return &dao.Invoice{
		ID:             i.ID,
		CompanyID:      i.CompanyID,
		ClientID:       i.ClientID,
		IssueDate:      i.IssueDate,
		PaymentAmount:  i.PaymentAmount,
		Fee:            i.Fee,
		FeeRate:        i.FeeRate,
		Tax:            i.Tax,
		TaxRate:        i.TaxRate,
		InvoiceAmount:  i.InvoiceAmount,
		PaymentDueDate: i.PaymentDueDate,
		Status:         i.Status,
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
	}
}

func InvoiceFromDAO(daoInvoice *dao.Invoice) *Invoice {
	return &Invoice{
		ID:             daoInvoice.ID,
		CompanyID:      daoInvoice.CompanyID,
		ClientID:       daoInvoice.ClientID,
		IssueDate:      daoInvoice.IssueDate,
		PaymentAmount:  daoInvoice.PaymentAmount,
		Fee:            daoInvoice.Fee,
		FeeRate:        daoInvoice.FeeRate,
		Tax:            daoInvoice.Tax,
		TaxRate:        daoInvoice.TaxRate,
		InvoiceAmount:  daoInvoice.InvoiceAmount,
		PaymentDueDate: daoInvoice.PaymentDueDate,
		Status:         daoInvoice.Status,
		CreatedAt:      daoInvoice.CreatedAt,
		UpdatedAt:      daoInvoice.UpdatedAt,
	}
}
