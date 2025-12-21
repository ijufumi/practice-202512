package gateway

import (
	"testing"

	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/dao"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupClientBankAccountTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:?_foreign_keys=on"), &gorm.Config{})
	assert.NoError(t, err)

	// マイグレーション
	err = db.AutoMigrate(&dao.Client{}, &dao.ClientBankAccount{})
	assert.NoError(t, err)

	return db
}

func TestClientBankAccountRepository_Create(t *testing.T) {
	db := setupClientBankAccountTestDB(t)
	repo := NewClientBankAccountRepository()

	// テストデータ準備
	company := &dao.Company{
		ID:            "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CorporateName: "Test Company",
	}
	err := db.Create(company).Error
	assert.NoError(t, err)
	client := &dao.Client{
		ID:                 "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CompanyID:          company.ID,
		CorporateName:      "Test Client",
		RepresentativeName: "Test Rep",
		PhoneNumber:        "000-0000-0000",
		PostalCode:         "000-0000",
		Address:            "Test Address",
	}
	err = db.Create(client).Error
	assert.NoError(t, err)

	t.Run("クライアント銀行口座作成成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		account := &models.ClientBankAccount{
			ClientID:      client.ID,
			BankName:      "Test Bank",
			BranchName:    "Test Branch",
			AccountNumber: "1234567890",
			AccountName:   "Test Account",
		}

		err := repo.Create(tx, account)
		assert.NoError(t, err)
		assert.NotEmpty(t, account.ID)
		assert.NotZero(t, account.CreatedAt)
		assert.NotZero(t, account.UpdatedAt)
	})
}

func TestClientBankAccountRepository_FindByID(t *testing.T) {
	db := setupClientBankAccountTestDB(t)
	repo := NewClientBankAccountRepository()

	// テストデータ準備
	company := &dao.Company{
		ID:            "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CorporateName: "Test Company",
	}
	err := db.Create(company).Error
	assert.NoError(t, err)
	client := &dao.Client{
		ID:                 "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CompanyID:          company.ID,
		CorporateName:      "Test Client",
		RepresentativeName: "Test Rep",
		PhoneNumber:        "000-0000-0000",
		PostalCode:         "000-0000",
		Address:            "Test Address",
	}
	err = db.Create(client).Error
	assert.NoError(t, err)

	testAccount := &dao.ClientBankAccount{
		ID:            "01HQZXFG0PJ9K8QXW7YM1N2ZXD",
		ClientID:      client.ID,
		BankName:      "Find Test Bank",
		BranchName:    "Find Test Branch",
		AccountNumber: "9876543210",
		AccountName:   "Find Test Account",
	}
	err = db.Create(testAccount).Error
	assert.NoError(t, err)

	t.Run("ID検索成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		account, err := repo.FindByID(tx, testAccount.ID)
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, testAccount.ID, account.ID)
		assert.Equal(t, testAccount.ClientID, account.ClientID)
		assert.Equal(t, testAccount.BankName, account.BankName)
		assert.Equal(t, testAccount.AccountNumber, account.AccountNumber)
	})

	t.Run("存在しないIDで失敗", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		account, err := repo.FindByID(tx, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
