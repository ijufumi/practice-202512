package gateway

import (
	"testing"

	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupCompanyTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:?_foreign_keys=on"), &gorm.Config{})
	assert.NoError(t, err)

	// マイグレーション
	err = db.AutoMigrate(&entities.Company{})
	assert.NoError(t, err)

	return db
}

func TestCompanyRepository_Create(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewCompanyRepository()

	t.Run("会社作成成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		company := &models.Company{
			CorporateName:      "Test Corporation",
			RepresentativeName: "Test Representative",
			PhoneNumber:        "000-0000-0000",
			PostalCode:         "000-0000",
			Address:            "Test Address",
		}

		err := repo.Create(tx, company)
		assert.NoError(t, err)
		assert.NotEmpty(t, company.ID)
		assert.NotZero(t, company.CreatedAt)
		assert.NotZero(t, company.UpdatedAt)
	})
}

func TestCompanyRepository_FindByID(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewCompanyRepository()

	// テストデータ準備
	testCompany := &entities.Company{
		ID:                 "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CorporateName:      "Find Test Corporation",
		RepresentativeName: "Find Test Representative",
		PhoneNumber:        "111-1111-1111",
		PostalCode:         "111-1111",
		Address:            "Find Test Address",
	}
	err := db.Create(testCompany).Error
	assert.NoError(t, err)

	t.Run("ID検索成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		company, err := repo.FindByID(tx, testCompany.ID)
		assert.NoError(t, err)
		assert.NotNil(t, company)
		assert.Equal(t, testCompany.ID, company.ID)
		assert.Equal(t, testCompany.CorporateName, company.CorporateName)
		assert.Equal(t, testCompany.RepresentativeName, company.RepresentativeName)
	})

	t.Run("存在しないIDで失敗", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		company, err := repo.FindByID(tx, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, company)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
