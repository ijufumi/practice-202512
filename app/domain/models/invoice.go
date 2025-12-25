package models

import (
	"github.com/ijufumi/practice-202512/app/domain/value"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/entities"
	"github.com/shopspring/decimal"

	"time"
)

type Invoice struct {
	ID             string
	CompanyID      string
	ClientID       string
	IssueDate      time.Time
	PaymentAmount  decimal.Decimal
	Fee            decimal.Decimal
	FeeRate        decimal.Decimal
	Tax            decimal.Decimal
	TaxRate        decimal.Decimal
	InvoiceAmount  decimal.Decimal
	PaymentDueDate time.Time
	Status         value.InvoiceStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (i *Invoice) ToDAO() *entities.Invoice {
	return &entities.Invoice{
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

func InvoiceFromDAO(daoInvoice *entities.Invoice) *Invoice {
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

// CalculateFee は支払金額に対する手数料を計算します
func (i *Invoice) CalculateFee(feeRate decimal.Decimal) {
	i.Fee = i.PaymentAmount.Mul(feeRate).Truncate(0)
	i.FeeRate = feeRate
}

// CalculateTax は手数料に対する消費税を計算します
func (i *Invoice) CalculateTax(taxRate decimal.Decimal) {
	i.Tax = i.Fee.Mul(taxRate).Truncate(0)
	i.TaxRate = taxRate
}

// CalculateInvoiceAmount は請求金額を計算します（支払金額 + 手数料 + 消費税）
func (i *Invoice) CalculateInvoiceAmount() {
	i.InvoiceAmount = i.PaymentAmount.Add(i.Fee).Add(i.Tax)
}
