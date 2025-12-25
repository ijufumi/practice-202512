package gateway

import (
	"testing"
	"time"

	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/domain/value"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/entities"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupInvoiceTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:?_foreign_keys=on"), &gorm.Config{})
	assert.NoError(t, err)

	// マイグレーション
	err = db.AutoMigrate(&entities.Client{}, &entities.Invoice{})
	assert.NoError(t, err)

	return db
}

func TestInvoiceRepository_Create(t *testing.T) {
	db := setupInvoiceTestDB(t)
	repo := NewInvoiceRepository()

	// テストデータ準備
	company := &entities.Company{
		ID:            "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CorporateName: "Test Company",
	}
	err := db.Create(company).Error
	assert.NoError(t, err)
	client := &entities.Client{
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

	t.Run("請求書作成成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		invoice := &models.Invoice{
			CompanyID:      company.ID,
			ClientID:       client.ID,
			IssueDate:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			PaymentAmount:  decimal.NewFromInt(100000),
			Fee:            decimal.NewFromInt(2000),
			FeeRate:        decimal.NewFromFloat(0.02),
			Tax:            decimal.NewFromInt(10000),
			TaxRate:        decimal.NewFromFloat(0.10),
			InvoiceAmount:  decimal.NewFromInt(112000),
			PaymentDueDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			Status:         value.InvoiceStatusUnprocessed,
		}

		err := repo.Create(tx, invoice)
		assert.NoError(t, err)
		assert.NotEmpty(t, invoice.ID)
		assert.NotZero(t, invoice.CreatedAt)
		assert.NotZero(t, invoice.UpdatedAt)
	})
}

func TestInvoiceRepository_FindByPaymentDueDateRange(t *testing.T) {
	db := setupInvoiceTestDB(t)
	repo := NewInvoiceRepository()

	// テストデータ準備
	company := &entities.Company{
		ID:            "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
		CorporateName: "Test Company",
	}
	err := db.Create(company).Error
	assert.NoError(t, err)
	client := &entities.Client{
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

	// 複数の請求書を作成
	invoices := []*entities.Invoice{
		{
			ID:             "01HQZXFG0PJ9K8QXW7YM1N2ZXD",
			CompanyID:      company.ID,
			ClientID:       client.ID,
			IssueDate:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			PaymentAmount:  decimal.NewFromInt(100000),
			Fee:            decimal.NewFromInt(2000),
			FeeRate:        decimal.NewFromFloat(0.02),
			Tax:            decimal.NewFromInt(10000),
			TaxRate:        decimal.NewFromFloat(0.10),
			InvoiceAmount:  decimal.NewFromInt(112000),
			PaymentDueDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			Status:         value.InvoiceStatusUnprocessed,
		},
		{
			ID:             "01HQZXFG0PJ9K8QXW7YM1N2ZXE",
			CompanyID:      company.ID,
			ClientID:       client.ID,
			IssueDate:      time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			PaymentAmount:  decimal.NewFromInt(200000),
			Fee:            decimal.NewFromInt(4000),
			FeeRate:        decimal.NewFromFloat(0.02),
			Tax:            decimal.NewFromInt(20000),
			TaxRate:        decimal.NewFromFloat(0.10),
			InvoiceAmount:  decimal.NewFromInt(224000),
			PaymentDueDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			Status:         value.InvoiceStatusUnprocessed,
		},
		{
			ID:             "01HQZXFG0PJ9K8QXW7YM1N2ZXF",
			CompanyID:      company.ID,
			ClientID:       client.ID,
			IssueDate:      time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			PaymentAmount:  decimal.NewFromInt(300000),
			Fee:            decimal.NewFromInt(6000),
			FeeRate:        decimal.NewFromFloat(0.02),
			Tax:            decimal.NewFromInt(30000),
			TaxRate:        decimal.NewFromFloat(0.10),
			InvoiceAmount:  decimal.NewFromInt(336000),
			PaymentDueDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
			Status:         value.InvoiceStatusUnprocessed,
		},
	}

	for _, inv := range invoices {
		err := db.Create(inv).Error
		assert.NoError(t, err)
	}

	var nilDate *time.Time

	t.Run("日付範囲で検索成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		startDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 3, 31, 23, 59, 59, 0, time.UTC)

		result, err := repo.FindByPaymentDueDateRange(tx, &startDate, &endDate, 0, 100)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, invoices[0].ID, result[0].ID)
		assert.Equal(t, invoices[1].ID, result[1].ID)
	})

	t.Run("開始日のみ指定", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		startDate := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)

		result, err := repo.FindByPaymentDueDateRange(tx, &startDate, nilDate, 0, 100)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, invoices[1].ID, result[0].ID)
		assert.Equal(t, invoices[2].ID, result[1].ID)
	})

	t.Run("終了日のみ指定", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		endDate := time.Date(2025, 2, 28, 23, 59, 59, 0, time.UTC)

		result, err := repo.FindByPaymentDueDateRange(tx, nilDate, &endDate, 0, 100)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, invoices[0].ID, result[0].ID)
	})

	t.Run("日付範囲指定なし（全件取得）", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		result, err := repo.FindByPaymentDueDateRange(tx, nilDate, nilDate, 0, 100)
		assert.NoError(t, err)
		assert.Len(t, result, 3)
	})

	t.Run("該当なし", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)

		result, err := repo.FindByPaymentDueDateRange(tx, &startDate, &endDate, 0, 100)
		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("offset/limitのテスト", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// offset=1, limit=1で2件目のみ取得
		result, err := repo.FindByPaymentDueDateRange(tx, nilDate, nilDate, 1, 1)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, invoices[1].ID, result[0].ID)
	})

	t.Run("limitで取得件数を制限", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// limit=2で最初の2件のみ取得
		result, err := repo.FindByPaymentDueDateRange(tx, nilDate, nilDate, 0, 2)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, invoices[0].ID, result[0].ID)
		assert.Equal(t, invoices[1].ID, result[1].ID)
	})
}
