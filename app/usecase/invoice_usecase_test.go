package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ijufumi/practice-202512/app/domain/models"
	repository "github.com/ijufumi/practice-202512/app/domain/repository/mocks"
	"github.com/ijufumi/practice-202512/app/domain/value"
	"github.com/ijufumi/practice-202512/app/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupInvoiceUsecaseContext(t *testing.T) (context.Context, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	ctx := util.SetDB(context.Background(), db)
	ctx = util.SetUserID(ctx, "userID")
	return ctx, db
}

func TestInvoiceUsecase_CreateInvoice(t *testing.T) {
	t.Run("請求書作成成功", func(t *testing.T) {
		ctx, _ := setupInvoiceUsecaseContext(t)
		mockInvoiceRepository := repository.NewMockInvoiceRepository(t)
		mockUserRepository := repository.NewMockUserRepository(t)

		clientID := "01HQZXFG0PJ9K8QXW7YM1N2ZXC"
		issueDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		paymentAmount := decimal.NewFromInt(100000)
		paymentDueDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)

		user := &models.User{
			ID:        "userID",
			CompanyID: "companyID",
		}
		mockUserRepository.EXPECT().FindByID(mock.Anything, user.ID).Return(user, nil)

		mockInvoiceRepository.EXPECT().Create(mock.Anything, mock.MatchedBy(func(inv *models.Invoice) bool {
			// 手数料と消費税の計算確認
			expectedFee := paymentAmount.Mul(decimal.NewFromFloat(0.04))             // 4000
			expectedTax := expectedFee.Mul(decimal.NewFromFloat(0.10))               // 400
			expectedInvoiceAmount := paymentAmount.Add(expectedFee).Add(expectedTax) // 104400

			return inv.ClientID == clientID &&
				inv.CompanyID == user.CompanyID &&
				inv.PaymentAmount.Equal(paymentAmount) &&
				inv.Fee.Equal(expectedFee) &&
				inv.Tax.Equal(expectedTax) &&
				inv.InvoiceAmount.Equal(expectedInvoiceAmount) &&
				inv.Status == value.InvoiceStatusUnprocessed
		})).Return(nil)

		usecase := NewInvoiceUsecase(mockInvoiceRepository, mockUserRepository)
		invoice, err := usecase.CreateInvoice(ctx, clientID, issueDate, paymentAmount, paymentDueDate)

		assert.NoError(t, err)
		assert.NotNil(t, invoice)
		assert.Equal(t, clientID, invoice.ClientID)
		assert.Equal(t, paymentAmount, invoice.PaymentAmount)
		assert.Equal(t, decimal.NewFromInt(4000), invoice.Fee)
		assert.True(t, invoice.FeeRate.Equal(decimal.NewFromFloat(0.04)))
		assert.Equal(t, decimal.NewFromInt(400), invoice.Tax)
		assert.True(t, invoice.TaxRate.Equal(decimal.NewFromFloat(0.10)))
		assert.Equal(t, decimal.NewFromInt(104400), invoice.InvoiceAmount)
		assert.Equal(t, value.InvoiceStatusUnprocessed, invoice.Status)
	})

	t.Run("手数料と消費税の計算確認 - 別の金額", func(t *testing.T) {
		ctx, _ := setupInvoiceUsecaseContext(t)
		mockInvoiceRepository := repository.NewMockInvoiceRepository(t)
		mockUserRepository := repository.NewMockUserRepository(t)

		clientID := "01HQZXFG0PJ9K8QXW7YM1N2ZXC"
		issueDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		paymentAmount := decimal.NewFromInt(250000) // 250,000円
		paymentDueDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)

		user := &models.User{
			ID:        "userID",
			CompanyID: "companyID",
		}
		mockUserRepository.EXPECT().FindByID(mock.Anything, user.ID).Return(user, nil)
		mockInvoiceRepository.EXPECT().Create(mock.Anything, mock.MatchedBy(func(inv *models.Invoice) bool {
			expectedFee := paymentAmount.Mul(decimal.NewFromFloat(0.04))             // 10000
			expectedTax := expectedFee.Mul(decimal.NewFromFloat(0.10))               // 1000
			expectedInvoiceAmount := paymentAmount.Add(expectedFee).Add(expectedTax) // 261000

			return inv.Fee.Equal(expectedFee) &&
				inv.Tax.Equal(expectedTax) &&
				inv.InvoiceAmount.Equal(expectedInvoiceAmount)
		})).Return(nil)

		usecase := NewInvoiceUsecase(mockInvoiceRepository, mockUserRepository)
		invoice, err := usecase.CreateInvoice(ctx, clientID, issueDate, paymentAmount, paymentDueDate)

		assert.NoError(t, err)
		assert.Equal(t, decimal.NewFromInt(10000), invoice.Fee)
		assert.Equal(t, decimal.NewFromInt(1000), invoice.Tax)
		assert.Equal(t, decimal.NewFromInt(261000), invoice.InvoiceAmount)
	})

	t.Run("リポジトリエラー", func(t *testing.T) {
		ctx, _ := setupInvoiceUsecaseContext(t)
		mockInvoiceRepository := repository.NewMockInvoiceRepository(t)
		mockUserRepository := repository.NewMockUserRepository(t)

		clientID := "01HQZXFG0PJ9K8QXW7YM1N2ZXC"
		issueDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		paymentAmount := decimal.NewFromInt(100000)
		paymentDueDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)

		user := &models.User{
			ID:        "userID",
			CompanyID: "companyID",
		}
		mockUserRepository.EXPECT().FindByID(mock.Anything, user.ID).Return(user, nil)

		mockInvoiceRepository.EXPECT().Create(mock.Anything, mock.Anything).
			Return(errors.New("database error"))

		usecase := NewInvoiceUsecase(mockInvoiceRepository, mockUserRepository)
		invoice, err := usecase.CreateInvoice(ctx, clientID, issueDate, paymentAmount, paymentDueDate)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Nil(t, invoice)
	})

	t.Run("コンテキストにDBがない", func(t *testing.T) {
		ctx := context.Background()
		mockInvoiceRepository := repository.NewMockInvoiceRepository(t)
		mockUserRepository := repository.NewMockUserRepository(t)

		clientID := "01HQZXFG0PJ9K8QXW7YM1N2ZXC"
		issueDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		paymentAmount := decimal.NewFromInt(100000)
		paymentDueDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)

		usecase := NewInvoiceUsecase(mockInvoiceRepository, mockUserRepository)
		invoice, err := usecase.CreateInvoice(ctx, clientID, issueDate, paymentAmount, paymentDueDate)

		assert.Error(t, err)
		assert.Equal(t, "database connection not found in context", err.Error())
		assert.Nil(t, invoice)
	})
}

func TestInvoiceUsecase_GetInvoicesByPaymentDueDateRange(t *testing.T) {
	var nilDate *time.Time

	t.Run("日付範囲で請求書取得成功", func(t *testing.T) {
		ctx, _ := setupInvoiceUsecaseContext(t)
		mockInvoiceRepository := repository.NewMockInvoiceRepository(t)
		mockUserRepository := repository.NewMockUserRepository(t)

		startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

		expectedInvoices := []*models.Invoice{
			{
				ID:             "01HQZXFG0PJ9K8QXW7YM1N2ZXD",
				ClientID:       "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
				PaymentAmount:  decimal.NewFromInt(100000),
				PaymentDueDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:             "01HQZXFG0PJ9K8QXW7YM1N2ZXE",
				ClientID:       "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
				PaymentAmount:  decimal.NewFromInt(200000),
				PaymentDueDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			},
		}

		mockInvoiceRepository.EXPECT().FindByPaymentDueDateRange(mock.Anything, &startDate, &endDate, 0, 100).
			Return(expectedInvoices, nil)

		usecase := NewInvoiceUsecase(mockInvoiceRepository, mockUserRepository)
		invoices, err := usecase.GetInvoicesByPaymentDueDateRange(ctx, &startDate, &endDate, 0, 100)

		assert.NoError(t, err)
		assert.Len(t, invoices, 2)
		assert.Equal(t, expectedInvoices[0].ID, invoices[0].ID)
		assert.Equal(t, expectedInvoices[1].ID, invoices[1].ID)
	})

	t.Run("日付範囲なし（全件取得）", func(t *testing.T) {
		ctx, _ := setupInvoiceUsecaseContext(t)
		mockInvoiceRepository := repository.NewMockInvoiceRepository(t)
		mockUserRepository := repository.NewMockUserRepository(t)

		expectedInvoices := []*models.Invoice{
			{ID: "01HQZXFG0PJ9K8QXW7YM1N2ZXD"},
			{ID: "01HQZXFG0PJ9K8QXW7YM1N2ZXE"},
			{ID: "01HQZXFG0PJ9K8QXW7YM1N2ZXF"},
		}

		mockInvoiceRepository.EXPECT().FindByPaymentDueDateRange(mock.Anything, nilDate, nilDate, 0, 100).
			Return(expectedInvoices, nil)

		usecase := NewInvoiceUsecase(mockInvoiceRepository, mockUserRepository)
		invoices, err := usecase.GetInvoicesByPaymentDueDateRange(ctx, nil, nil, 0, 100)

		assert.NoError(t, err)
		assert.Len(t, invoices, 3)
	})

	t.Run("リポジトリエラー", func(t *testing.T) {
		ctx, _ := setupInvoiceUsecaseContext(t)
		mockInvoiceRepository := repository.NewMockInvoiceRepository(t)
		mockUserRepository := repository.NewMockUserRepository(t)

		startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

		mockInvoiceRepository.EXPECT().FindByPaymentDueDateRange(mock.Anything, &startDate, &endDate, 0, 100).
			Return(nil, errors.New("database error"))

		usecase := NewInvoiceUsecase(mockInvoiceRepository, mockUserRepository)
		invoices, err := usecase.GetInvoicesByPaymentDueDateRange(ctx, &startDate, &endDate, 0, 100)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Nil(t, invoices)
	})

	t.Run("コンテキストにDBがない", func(t *testing.T) {
		ctx := context.Background()
		mockInvoiceRepository := repository.NewMockInvoiceRepository(t)
		mockUserRepository := repository.NewMockUserRepository(t)

		usecase := NewInvoiceUsecase(mockInvoiceRepository, mockUserRepository)
		invoices, err := usecase.GetInvoicesByPaymentDueDateRange(ctx, nil, nil, 0, 100)

		assert.Error(t, err)
		assert.Equal(t, "database connection not found in context", err.Error())
		assert.Nil(t, invoices)
	})

	t.Run("offsetとlimitのデフォルト値設定", func(t *testing.T) {
		ctx, _ := setupInvoiceUsecaseContext(t)
		mockInvoiceRepository := repository.NewMockInvoiceRepository(t)
		mockUserRepository := repository.NewMockUserRepository(t)

		expectedInvoices := []*models.Invoice{
			{ID: "01HQZXFG0PJ9K8QXW7YM1N2ZXD"},
		}

		// 負のoffsetは0に、0以下のlimitは100に補正される
		mockInvoiceRepository.EXPECT().FindByPaymentDueDateRange(mock.Anything, nilDate, nilDate, 0, 100).
			Return(expectedInvoices, nil)

		usecase := NewInvoiceUsecase(mockInvoiceRepository, mockUserRepository)
		invoices, err := usecase.GetInvoicesByPaymentDueDateRange(ctx, nil, nil, -1, 0)

		assert.NoError(t, err)
		assert.Len(t, invoices, 1)
	})

	t.Run("カスタムoffsetとlimit", func(t *testing.T) {
		ctx, _ := setupInvoiceUsecaseContext(t)
		mockInvoiceRepository := repository.NewMockInvoiceRepository(t)
		mockUserRepository := repository.NewMockUserRepository(t)

		expectedInvoices := []*models.Invoice{
			{ID: "01HQZXFG0PJ9K8QXW7YM1N2ZXE"},
		}

		mockInvoiceRepository.EXPECT().FindByPaymentDueDateRange(mock.Anything, nilDate, nilDate, 10, 20).
			Return(expectedInvoices, nil)

		usecase := NewInvoiceUsecase(mockInvoiceRepository, mockUserRepository)
		invoices, err := usecase.GetInvoicesByPaymentDueDateRange(ctx, nil, nil, 10, 20)

		assert.NoError(t, err)
		assert.Len(t, invoices, 1)
	})
}
