package entities

import (
	"time"

	"github.com/ijufumi/practice-202512/app/util"

	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey;type:char(26)" json:"id"`
	CompanyID string    `gorm:"type:char(26);not null;index" json:"company_id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Company Company `gorm:"foreignKey:CompanyID"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = util.GenerateULID()
	}

	return nil
}
