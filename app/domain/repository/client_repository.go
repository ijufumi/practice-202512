package repository

import (
	"github.com/ijufumi/practice-202512/app/domain/models"

	"gorm.io/gorm"
)

type ClientRepository interface {
	Create(db *gorm.DB, client *models.Client) error
	FindByID(db *gorm.DB, id string) (*models.Client, error)
}
