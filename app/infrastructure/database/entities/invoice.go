package entities

import (
	"time"

	"github.com/ijufumi/practice-202512/app/domain/value"
	"github.com/ijufumi/practice-202512/app/util"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

type Invoice struct {
	ID             string              `gorm:"primaryKey;type:char(26)" json:"id"`
	CompanyID      string              `gorm:"type:char(26);not null;index" json:"company_id"`
	ClientID       string              `gorm:"type:char(26);not null;index" json:"client_id"`
	IssueDate      time.Time           `gorm:"not null" json:"issue_date"`
	PaymentAmount  decimal.Decimal     `gorm:"type:decimal(20,2);not null" json:"payment_amount"`
	Fee            decimal.Decimal     `gorm:"type:decimal(20,2);not null" json:"fee"`
	FeeRate        decimal.Decimal     `gorm:"type:decimal(5,4);not null" json:"fee_rate"`
	Tax            decimal.Decimal     `gorm:"type:decimal(20,2);not null" json:"tax"`
	TaxRate        decimal.Decimal     `gorm:"type:decimal(5,4);not null" json:"tax_rate"`
	InvoiceAmount  decimal.Decimal     `gorm:"type:decimal(20,2);not null" json:"invoice_amount"`
	PaymentDueDate time.Time           `gorm:"not null;index" json:"payment_due_date"`
	Status         value.InvoiceStatus `gorm:"size:20;not null;index" json:"status"`
	CreatedAt      time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time           `gorm:"autoUpdateTime" json:"updated_at"`

	Company Company `gorm:"foreignKey:CompanyID"`
	Client  Client  `gorm:"foreignKey:ClientID"`
}

func (i *Invoice) TableName() string {
	return "invoices"
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = util.GenerateULID()
	}
	return nil
}
