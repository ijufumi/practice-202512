package gateway

import (
	"testing"

	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/dao"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:?_foreign_keys=on"), &gorm.Config{})
	assert.NoError(t, err)

	// マイグレーション
	err = db.AutoMigrate(&dao.Company{}, &dao.User{})
	assert.NoError(t, err)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository()

	// テストデータ準備
	company := &dao.Company{
		ID:            "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CorporateName: "Test Company",
	}
	err := db.Create(company).Error
	assert.NoError(t, err)

	t.Run("ユーザー作成成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		user := &models.User{
			CompanyID: company.ID,
			Name:      "Test User",
			Email:     "test@example.com",
			Password:  "hashedpassword",
		}

		err := repo.Create(tx, user)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)
	})

	t.Run("重複するEmailで失敗", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		user1 := &models.User{
			CompanyID: company.ID,
			Name:      "User 1",
			Email:     "duplicate@example.com",
			Password:  "password1",
		}
		err := repo.Create(tx, user1)
		assert.NoError(t, err)

		user2 := &models.User{
			CompanyID: company.ID,
			Name:      "User 2",
			Email:     "duplicate@example.com",
			Password:  "password2",
		}
		err = repo.Create(tx, user2)
		assert.Error(t, err)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository()

	// テストデータ準備
	company := &dao.Company{
		ID:            "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CorporateName: "Test Company",
	}
	err := db.Create(company).Error
	assert.NoError(t, err)

	testUser := &dao.User{
		ID:        "01HQZXFG0PJ9K8QXW7YM1N2ZXD",
		CompanyID: company.ID,
		Name:      "Find Test User",
		Email:     "findtest@example.com",
		Password:  "password",
	}
	err = db.Create(testUser).Error
	assert.NoError(t, err)

	t.Run("ID検索成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		user, err := repo.FindByID(tx, testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Email, user.Email)
	})

	t.Run("存在しないIDで失敗", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		user, err := repo.FindByID(tx, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository()

	// テストデータ準備
	company := &dao.Company{
		ID:            "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CorporateName: "Test Company",
	}
	err := db.Create(company).Error
	assert.NoError(t, err)

	testUser := &dao.User{
		ID:        "01HQZXFG0PJ9K8QXW7YM1N2ZXE",
		CompanyID: company.ID,
		Name:      "Email Test User",
		Email:     "emailtest@example.com",
		Password:  "password",
	}
	err = db.Create(testUser).Error
	assert.NoError(t, err)

	t.Run("Email検索成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		user, err := repo.FindByEmail(tx, testUser.Email)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUser.Email, user.Email)
		assert.Equal(t, testUser.ID, user.ID)
	})

	t.Run("存在しないEmailで失敗", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		user, err := repo.FindByEmail(tx, "nonexistent@example.com")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
