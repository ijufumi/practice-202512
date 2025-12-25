package gateway

import (
	"testing"

	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupClientTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:?_foreign_keys=on"), &gorm.Config{})
	assert.NoError(t, err)

	// マイグレーション
	err = db.AutoMigrate(&entities.Client{})
	assert.NoError(t, err)

	return db
}

func TestClientRepository_Create(t *testing.T) {
	db := setupClientTestDB(t)
	repo := NewClientRepository()

	company := &entities.Company{
		ID:            "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CorporateName: "Test Company",
	}
	err := db.Create(company).Error
	assert.NoError(t, err)

	t.Run("クライアント作成成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		client := &models.Client{
			CompanyID:          company.ID,
			CorporateName:      "Test Corporation",
			RepresentativeName: "Test Representative",
			PhoneNumber:        "000-0000-0000",
			PostalCode:         "000-0000",
			Address:            "Test Address",
		}

		err := repo.Create(tx, client)
		assert.NoError(t, err)
		assert.NotEmpty(t, client.ID)
		assert.NotZero(t, client.CreatedAt)
		assert.NotZero(t, client.UpdatedAt)
	})
}

func TestClientRepository_FindByID(t *testing.T) {
	db := setupClientTestDB(t)
	repo := NewClientRepository()

	// テストデータ準備
	company := &entities.Company{
		ID:            "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CorporateName: "Test Company",
	}
	err := db.Create(company).Error
	assert.NoError(t, err)
	testClient := &entities.Client{
		ID:                 "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CompanyID:          company.ID,
		CorporateName:      "Find Test Corporation",
		RepresentativeName: "Find Test Representative",
		PhoneNumber:        "111-1111-1111",
		PostalCode:         "111-1111",
		Address:            "Find Test Address",
	}
	err = db.Create(testClient).Error
	assert.NoError(t, err)

	t.Run("ID検索成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		client, err := repo.FindByID(tx, testClient.ID)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, testClient.ID, client.ID)
		assert.Equal(t, testClient.CorporateName, client.CorporateName)
		assert.Equal(t, testClient.RepresentativeName, client.RepresentativeName)
	})

	t.Run("存在しないIDで失敗", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		client, err := repo.FindByID(tx, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
