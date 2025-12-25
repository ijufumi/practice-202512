package models

import (
	domainModel "github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/domain/value"
	"github.com/shopspring/decimal"

	"time"
)

type CreateInvoiceRequest struct {
	ClientID       string          `json:"client_id" validate:"required"`
	IssueDate      string          `json:"issue_date" validate:"required"`
	PaymentAmount  decimal.Decimal `json:"payment_amount" validate:"required,min=1"`
	PaymentDueDate string          `json:"payment_due_date" validate:"required"`
}

type InvoiceResponse struct {
	ID             string              `json:"id"`
	ClientID       string              `json:"client_id"`
	IssueDate      time.Time           `json:"issue_date"`
	PaymentAmount  decimal.Decimal     `json:"payment_amount"`
	Fee            decimal.Decimal     `json:"fee"`
	FeeRate        decimal.Decimal     `json:"fee_rate"`
	Tax            decimal.Decimal     `json:"tax"`
	TaxRate        decimal.Decimal     `json:"tax_rate"`
	InvoiceAmount  decimal.Decimal     `json:"invoice_amount"`
	PaymentDueDate time.Time           `json:"payment_due_date"`
	Status         value.InvoiceStatus `json:"status"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

func FromInvoiceDomainModel(invoice *domainModel.Invoice) *InvoiceResponse {
	return &InvoiceResponse{
		ID:             invoice.ID,
		ClientID:       invoice.ClientID,
		IssueDate:      invoice.IssueDate,
		PaymentAmount:  invoice.PaymentAmount,
		Fee:            invoice.Fee,
		FeeRate:        invoice.FeeRate,
		Tax:            invoice.Tax,
		TaxRate:        invoice.TaxRate,
		InvoiceAmount:  invoice.InvoiceAmount,
		PaymentDueDate: invoice.PaymentDueDate,
		Status:         invoice.Status,
		CreatedAt:      invoice.CreatedAt,
		UpdatedAt:      invoice.UpdatedAt,
	}
}

func FromInvoiceDomainModels(invoices []*domainModel.Invoice) []*InvoiceResponse {
	responses := make([]*InvoiceResponse, len(invoices))
	for i, invoice := range invoices {
		responses[i] = FromInvoiceDomainModel(invoice)
	}
	return responses
}
