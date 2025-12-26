package entities

import (
	"time"

	"github.com/ijufumi/practice-202512/app/util"

	"gorm.io/gorm"
)

type ClientBankAccount struct {
	ID            string    `gorm:"primaryKey;type:char(26)" json:"id"`
	ClientID      string    `gorm:"type:char(26);not null;index" json:"client_id"`
	BankName      string    `gorm:"size:100;not null" json:"bank_name"`
	BranchName    string    `gorm:"size:100;not null" json:"branch_name"`
	AccountNumber string    `gorm:"size:20;not null" json:"account_number"`
	AccountName   string    `gorm:"size:100;not null" json:"account_name"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Client Client `gorm:"foreignKey:ClientID"`
}

func (c *ClientBankAccount) TableName() string {
	return "client_bank_accounts"
}

func (c *ClientBankAccount) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = util.GenerateULID()
	}

	return nil
}
