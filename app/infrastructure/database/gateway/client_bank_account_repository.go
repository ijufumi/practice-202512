package gateway

import (
	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/domain/repository"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/entities"

	"gorm.io/gorm"
)

type clientBankAccountRepository struct{}

func NewClientBankAccountRepository() repository.ClientBankAccountRepository {
	return &clientBankAccountRepository{}
}

func (r *clientBankAccountRepository) Create(db *gorm.DB, account *models.ClientBankAccount) error {
	daoAccount := account.ToDAO()
	if err := db.Create(&daoAccount).Error; err != nil {
		return err
	}
	account.ID = daoAccount.ID
	account.CreatedAt = daoAccount.CreatedAt
	account.UpdatedAt = daoAccount.UpdatedAt

	return nil
}

func (r *clientBankAccountRepository) FindByID(db *gorm.DB, id string) (*models.ClientBankAccount, error) {
	var daoAccount entities.ClientBankAccount
	if err := db.First(&daoAccount, "id = ?", id).Error; err != nil {
		return nil, err
	}
	account := models.ClientBankAccountFromDAO(&daoAccount)

	return account, nil
}
