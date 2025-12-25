package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/domain/value"
	usecase "github.com/ijufumi/practice-202512/app/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInvoiceHandler_CreateInvoice(t *testing.T) {
	t.Run("請求書作成成功", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		clientID := "01HQZXFG0PJ9K8QXW7YM1N2ZXC"
		issueDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		paymentAmount := decimal.NewFromInt(100000)
		paymentDueDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)

		expectedInvoice := &models.Invoice{
			ID:             "01HQZXFG0PJ9K8QXW7YM1N2ZXD",
			ClientID:       clientID,
			IssueDate:      issueDate,
			PaymentAmount:  paymentAmount,
			Fee:            decimal.NewFromInt(4000),
			FeeRate:        decimal.NewFromFloat(0.04),
			Tax:            decimal.NewFromInt(400),
			TaxRate:        decimal.NewFromFloat(0.10),
			InvoiceAmount:  decimal.NewFromInt(104400),
			PaymentDueDate: paymentDueDate,
			Status:         value.InvoiceStatusUnprocessed,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockUsecase.EXPECT().CreateInvoice(
			mock.Anything,
			clientID,
			issueDate,
			paymentAmount,
			paymentDueDate,
		).Return(expectedInvoice, nil)

		handler := NewInvoiceHandler(mockUsecase)

		reqBody := `{
			"client_id": "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
			"issue_date": "2025-01-01",
			"payment_amount": 100000,
			"payment_due_date": "2025-02-01"
		}`
		req := httptest.NewRequest(http.MethodPost, "/invoices", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateInvoice(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedInvoice.ID, response["id"])
		assert.Equal(t, expectedInvoice.PaymentAmount.String(), response["payment_amount"])
		assert.Equal(t, expectedInvoice.InvoiceAmount.String(), response["invoice_amount"])
	})

	t.Run("不正なリクエストボディ", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		handler := NewInvoiceHandler(mockUsecase)

		reqBody := `{"invalid": "json"}`
		req := httptest.NewRequest(http.MethodPost, "/invoices", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateInvoice(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("バリデーションエラー - 必須フィールド不足", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		handler := NewInvoiceHandler(mockUsecase)

		reqBody := `{
			"payment_amount": 100000
		}`
		req := httptest.NewRequest(http.MethodPost, "/invoices", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateInvoice(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("バリデーションエラー - 不正な金額", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		handler := NewInvoiceHandler(mockUsecase)

		reqBody := `{
			"client_id": "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
			"issue_date": "2025-01-01",
			"payment_amount": "*",
			"payment_due_date": "2025-02-01"
		}`
		req := httptest.NewRequest(http.MethodPost, "/invoices", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateInvoice(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("不正な日付フォーマット - issue_date", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		handler := NewInvoiceHandler(mockUsecase)

		reqBody := `{
			"client_id": "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
			"issue_date": "2025/01/01",
			"payment_amount": 100000,
			"payment_due_date": "2025-02-01"
		}`
		req := httptest.NewRequest(http.MethodPost, "/invoices", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateInvoice(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid issue_date format")
	})

	t.Run("不正な日付フォーマット - payment_due_date", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		handler := NewInvoiceHandler(mockUsecase)

		reqBody := `{
			"client_id": "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
			"issue_date": "2025-01-01",
			"payment_amount": 100000,
			"payment_due_date": "2025/02/01"
		}`
		req := httptest.NewRequest(http.MethodPost, "/invoices", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateInvoice(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid payment_due_date format")
	})

	t.Run("Usecaseエラー", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		mockUsecase.EXPECT().CreateInvoice(
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, errors.New("database error"))

		handler := NewInvoiceHandler(mockUsecase)

		reqBody := `{
			"client_id": "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
			"issue_date": "2025-01-01",
			"payment_amount": 100000,
			"payment_due_date": "2025-02-01"
		}`
		req := httptest.NewRequest(http.MethodPost, "/invoices", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateInvoice(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestInvoiceHandler_GetInvoices(t *testing.T) {
	t.Run("請求書一覧取得成功 - 日付範囲指定あり", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

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

		mockUsecase.EXPECT().GetInvoicesByPaymentDueDateRange(
			mock.Anything,
			&startDate,
			&endDate,
			0,
			100,
		).Return(expectedInvoices, nil)

		handler := NewInvoiceHandler(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/invoices?start_date=2025-01-01&end_date=2025-12-31", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetInvoices(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, expectedInvoices[0].ID, response[0]["id"])
	})

	t.Run("請求書一覧取得成功 - 日付範囲指定なし", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		expectedInvoices := []*models.Invoice{
			{ID: "01HQZXFG0PJ9K8QXW7YM1N2ZXD"},
			{ID: "01HQZXFG0PJ9K8QXW7YM1N2ZXE"},
		}

		var nilTime *time.Time
		mockUsecase.EXPECT().GetInvoicesByPaymentDueDateRange(
			mock.Anything,
			nilTime,
			nilTime,
			0,
			100,
		).Return(expectedInvoices, nil)

		handler := NewInvoiceHandler(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/invoices", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetInvoices(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
	})

	t.Run("不正な日付フォーマット - start_date", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		handler := NewInvoiceHandler(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/invoices?start_date=2025/01/01", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetInvoices(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid start_date format")
	})

	t.Run("不正な日付フォーマット - end_date", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		handler := NewInvoiceHandler(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/invoices?end_date=2025/12/31", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetInvoices(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid end_date format")
	})

	t.Run("Usecaseエラー", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		mockUsecase.EXPECT().GetInvoicesByPaymentDueDateRange(
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, errors.New("database error"))

		handler := NewInvoiceHandler(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/invoices", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetInvoices(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("offsetとlimitのパラメータ指定", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		var nilTime *time.Time
		expectedInvoices := []*models.Invoice{
			{ID: "01HQZXFG0PJ9K8QXW7YM1N2ZXD"},
		}

		mockUsecase.EXPECT().GetInvoicesByPaymentDueDateRange(
			mock.Anything,
			nilTime,
			nilTime,
			10,
			20,
		).Return(expectedInvoices, nil)

		handler := NewInvoiceHandler(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/invoices?offset=10&limit=20", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetInvoices(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("不正なoffsetパラメータ", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		handler := NewInvoiceHandler(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/invoices?offset=invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetInvoices(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid offset parameter")
	})

	t.Run("不正なlimitパラメータ", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		handler := NewInvoiceHandler(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/invoices?limit=invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetInvoices(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid limit parameter")
	})

	t.Run("負のoffsetパラメータ", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		handler := NewInvoiceHandler(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/invoices?offset=-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetInvoices(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid offset parameter")
	})

	t.Run("0以下のlimitパラメータ", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockInvoiceUsecase(t)

		handler := NewInvoiceHandler(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/invoices?limit=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetInvoices(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid limit parameter")
	})
}
