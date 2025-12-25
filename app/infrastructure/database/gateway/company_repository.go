package gateway

import (
	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/domain/repository"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/entities"

	"gorm.io/gorm"
)

type companyRepository struct{}

func NewCompanyRepository() repository.CompanyRepository {
	return &companyRepository{}
}

func (r *companyRepository) Create(db *gorm.DB, company *models.Company) error {
	daoCompany := company.ToDAO()
	if err := db.Create(&daoCompany).Error; err != nil {
		return err
	}
	company.ID = daoCompany.ID
	company.CreatedAt = daoCompany.CreatedAt
	company.UpdatedAt = daoCompany.UpdatedAt
	return nil
}

func (r *companyRepository) FindByID(db *gorm.DB, id string) (*models.Company, error) {
	var daoCompany entities.Company
	if err := db.First(&daoCompany, "id = ?", id).Error; err != nil {
		return nil, err
	}
	company := models.CompanyFromDAO(&daoCompany)
	return company, nil
}
