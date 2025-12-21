package gateway

import (
	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/domain/repository"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/dao"

	"gorm.io/gorm"
)

type clientRepository struct{}

func NewClientRepository() repository.ClientRepository {
	return &clientRepository{}
}

func (r *clientRepository) Create(db *gorm.DB, client *models.Client) error {
	daoClient := client.ToDAO()
	if err := db.Create(&daoClient).Error; err != nil {
		return err
	}
	client.ID = daoClient.ID
	client.CreatedAt = daoClient.CreatedAt
	client.UpdatedAt = daoClient.UpdatedAt
	return nil
}

func (r *clientRepository) FindByID(db *gorm.DB, id string) (*models.Client, error) {
	var daoClient dao.Client
	if err := db.First(&daoClient, "id = ?", id).Error; err != nil {
		return nil, err
	}
	client := models.ClientFromDAO(&daoClient)
	return client, nil
}
