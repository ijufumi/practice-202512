package repository

import (
	"github.com/ijufumi/practice-202512/app/domain/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(db *gorm.DB, user *models.User) error
	FindByID(db *gorm.DB, id string) (*models.User, error)
	FindByEmail(db *gorm.DB, email string) (*models.User, error)
}
