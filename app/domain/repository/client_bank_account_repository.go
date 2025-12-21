package repository

import (
	"github.com/ijufumi/practice-202512/app/domain/models"

	"gorm.io/gorm"
)

type ClientBankAccountRepository interface {
	Create(db *gorm.DB, account *models.ClientBankAccount) error
	FindByID(db *gorm.DB, id string) (*models.ClientBankAccount, error)
}
