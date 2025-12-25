package gateway

import (
	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/domain/repository"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/entities"

	"gorm.io/gorm"
)

type userRepository struct{}

func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(db *gorm.DB, user *models.User) error {
	daoUser := user.ToDAO()
	if err := db.Create(&daoUser).Error; err != nil {
		return err
	}
	user.ID = daoUser.ID
	user.CreatedAt = daoUser.CreatedAt
	user.UpdatedAt = daoUser.UpdatedAt
	return nil
}

func (r *userRepository) FindByID(db *gorm.DB, id string) (*models.User, error) {
	var daoUser entities.User
	if err := db.First(&daoUser, "id = ?", id).Error; err != nil {
		return nil, err
	}
	user := models.UserFromDAO(&daoUser)
	return user, nil
}

func (r *userRepository) FindByEmail(db *gorm.DB, email string) (*models.User, error) {
	var daoUser entities.User
	if err := db.Where("email = ?", email).First(&daoUser).Error; err != nil {
		return nil, err
	}
	user := models.UserFromDAO(&daoUser)
	return user, nil
}
