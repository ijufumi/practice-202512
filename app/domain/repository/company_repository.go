package repository

import (
	"github.com/ijufumi/practice-202512/app/domain/models"

	"gorm.io/gorm"
)

type CompanyRepository interface {
	Create(db *gorm.DB, company *models.Company) error
	FindByID(db *gorm.DB, id string) (*models.Company, error)
}
