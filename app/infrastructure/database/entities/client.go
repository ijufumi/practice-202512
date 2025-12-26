package entities

import (
	"time"

	"github.com/ijufumi/practice-202512/app/util"

	"gorm.io/gorm"
)

type Client struct {
	ID                 string    `gorm:"primaryKey;type:char(26)" json:"id"`
	CompanyID          string    `gorm:"type:char(26);not null;index" json:"company_id"`
	CorporateName      string    `gorm:"size:200;not null" json:"corporate_name"`
	RepresentativeName string    `gorm:"size:100;not null" json:"representative_name"`
	PhoneNumber        string    `gorm:"size:20;not null" json:"phone_number"`
	PostalCode         string    `gorm:"size:10;not null" json:"postal_code"`
	Address            string    `gorm:"size:500;not null" json:"address"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Company Company `gorm:"foreignKey:CompanyID"`
}

func (c *Client) TableName() string {
	return "clients"
}

func (c *Client) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = util.GenerateULID()
	}

	return nil
}
