package dao

import (
	"time"

	"github.com/ijufumi/practice-202512/app/domain/value"
	"github.com/ijufumi/practice-202512/app/util"

	"gorm.io/gorm"
)

type Invoice struct {
	ID             string              `gorm:"primaryKey;type:char(26)" json:"id"`
	CompanyID      string              `gorm:"type:char(26);not null;index" json:"company_id"`
	ClientID       string              `gorm:"type:char(26);not null;index" json:"client_id"`
	IssueDate      time.Time           `gorm:"not null" json:"issue_date"`
	PaymentAmount  int                 `gorm:"not null" json:"payment_amount"`
	Fee            int                 `gorm:"not null" json:"fee"`
	FeeRate        float64             `gorm:"type:decimal(5,4);not null" json:"fee_rate"`
	Tax            int                 `gorm:"not null" json:"tax"`
	TaxRate        float64             `gorm:"type:decimal(5,4);not null" json:"tax_rate"`
	InvoiceAmount  int                 `gorm:"not null" json:"invoice_amount"`
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
