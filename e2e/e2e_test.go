package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ijufumi/practice-202512/app/config"
	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/gateway"
	"github.com/ijufumi/practice-202512/app/presentation"
	"github.com/ijufumi/practice-202512/app/presentation/handler"
	"github.com/ijufumi/practice-202512/app/usecase"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ijufumi/practice-202512/app/infrastructure/database/dao"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:?_foreign_keys=on"), &gorm.Config{})
	assert.NoError(t, err)

	// マイグレーション
	err = db.AutoMigrate(
		&dao.Company{},
		&dao.User{},
		&dao.Client{},
		&dao.ClientBankAccount{},
		&dao.Invoice{},
	)
	assert.NoError(t, err)

	return db
}

func setupTestData(t *testing.T, db *gorm.DB) (string, string) {
	// 会社データを作成
	company := &models.Company{
		CorporateName:      "Test Corporation",
		RepresentativeName: "Test Representative",
		PhoneNumber:        "000-0000-0000",
		PostalCode:         "000-0000",
		Address:            "Test Address",
	}
	companyRepo := gateway.NewCompanyRepository()
	err := companyRepo.Create(db, company)
	assert.NoError(t, err)

	// ユーザーデータを作成
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	user := &models.User{
		CompanyID: company.ID,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  string(passwordHash),
	}
	userRepo := gateway.NewUserRepository()
	err = userRepo.Create(db, user)
	assert.NoError(t, err)

	// クライアントデータを作成
	client := &models.Client{
		CompanyID:          company.ID,
		CorporateName:      "Client Corporation",
		RepresentativeName: "Client Representative",
		PhoneNumber:        "111-1111-1111",
		PostalCode:         "111-1111",
		Address:            "Client Address",
	}
	clientRepo := gateway.NewClientRepository()
	err = clientRepo.Create(db, client)
	assert.NoError(t, err)

	// クライアント銀行口座データを作成
	clientBankAccount := &models.ClientBankAccount{
		ClientID:      client.ID,
		BankName:      "Test Bank",
		BranchName:    "Test Branch",
		AccountNumber: "1234567890",
		AccountName:   "Test Account",
	}
	clientBankAccountRepo := gateway.NewClientBankAccountRepository()
	err = clientBankAccountRepo.Create(db, clientBankAccount)
	assert.NoError(t, err)

	return user.Email, client.ID
}

func setupRouter(db *gorm.DB, cfg *config.Config) *httptest.Server {
	// 依存性の注入
	invoiceRepository := gateway.NewInvoiceRepository()
	userRepository := gateway.NewUserRepository()
	invoiceUsecase := usecase.NewInvoiceUsecase(invoiceRepository, userRepository)
	invoiceHandler := handler.NewInvoiceHandler(invoiceUsecase)

	authUsecase := usecase.NewAuthUsecase(userRepository, cfg)
	authHandler := handler.NewAuthHandler(authUsecase)

	router := presentation.NewRouter(db, cfg, invoiceHandler, authHandler)

	return httptest.NewServer(router)
}

func TestE2E_LoginAndCreateInvoice(t *testing.T) {
	// テスト用DBのセットアップ
	db := setupTestDB(t)

	// テストデータの作成
	email, clientID := setupTestData(t, db)

	// テスト用の設定
	cfg := &config.Config{
		JWTSecret: "test-secret-key-for-e2e",
	}

	// サーバーのセットアップ
	server := setupRouter(db, cfg)
	defer server.Close()

	t.Run("E2E - ログインから請求書作成まで", func(t *testing.T) {
		// Step 1: ログイン
		loginReq := map[string]string{
			"email":    email,
			"password": "testpassword",
		}
		loginBody, _ := json.Marshal(loginReq)

		resp, err := http.Post(
			server.URL+"/api/login",
			"application/json",
			bytes.NewBuffer(loginBody),
		)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResp map[string]string
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		assert.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		token := loginResp["token"]
		assert.NotEmpty(t, token)

		// Step 2: 請求書作成
		invoiceReq := map[string]interface{}{
			"client_id":        clientID,
			"issue_date":       time.Now().Format(time.DateOnly),
			"payment_amount":   100000,
			"payment_due_date": time.Now().AddDate(0, 1, 0).Format(time.DateOnly),
		}
		invoiceBody, _ := json.Marshal(invoiceReq)

		req, _ := http.NewRequest(
			http.MethodPost,
			server.URL+"/api/invoices",
			bytes.NewBuffer(invoiceBody),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var invoiceResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&invoiceResp)
		assert.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.NotEmpty(t, invoiceResp["id"])
		assert.Equal(t, float64(100000), invoiceResp["payment_amount"])

		// Step 3: 請求書一覧取得
		req, _ = http.NewRequest(
			http.MethodGet,
			server.URL+"/api/invoices",
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var invoicesResp []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&invoicesResp)
		assert.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.NotEmpty(t, invoicesResp)
		assert.Equal(t, invoiceResp["id"], invoicesResp[0]["id"])
	})

	t.Run("E2E - offsetとlimitでページネーション", func(t *testing.T) {
		// Step 1: ログイン
		loginReq := map[string]string{
			"email":    email,
			"password": "testpassword",
		}
		loginBody, _ := json.Marshal(loginReq)

		resp, err := http.Post(
			server.URL+"/api/login",
			"application/json",
			bytes.NewBuffer(loginBody),
		)
		assert.NoError(t, err)
		var loginResp map[string]string
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		assert.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		token := loginResp["token"]

		// Step 2: 複数の請求書を作成
		for i := 0; i < 3; i++ {
			invoiceReq := map[string]interface{}{
				"client_id":        clientID,
				"issue_date":       time.Now().Format(time.DateOnly),
				"payment_amount":   100000 * (i + 1),
				"payment_due_date": time.Now().AddDate(0, i+1, 0).Format(time.DateOnly),
			}
			invoiceBody, _ := json.Marshal(invoiceReq)

			req, _ := http.NewRequest(
				http.MethodPost,
				server.URL+"/api/invoices",
				bytes.NewBuffer(invoiceBody),
			)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			client := &http.Client{}
			resp, _ := client.Do(req)
			defer func() { _ = resp.Body.Close() }()
		}

		// Step 3: offset=1, limit=1で取得
		req, _ := http.NewRequest(
			http.MethodGet,
			server.URL+"/api/invoices?offset=1&limit=1",
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var invoicesResp []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&invoicesResp)
		assert.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		// 1件のみ取得されることを確認
		assert.Len(t, invoicesResp, 1)
	})

	t.Run("E2E - 認証エラー", func(t *testing.T) {
		// 無効なログイン
		loginReq := map[string]string{
			"email":    email,
			"password": "wrongpassword",
		}
		loginBody, _ := json.Marshal(loginReq)

		resp, err := http.Post(
			server.URL+"/api/login",
			"application/json",
			bytes.NewBuffer(loginBody),
		)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		defer func() { _ = resp.Body.Close() }()
	})

	t.Run("E2E - JWT認証なしでAPI呼び出し", func(t *testing.T) {
		// JWTトークンなしで請求書作成
		invoiceReq := map[string]interface{}{
			"client_id":        clientID,
			"issue_date":       time.Now().Format(time.DateOnly),
			"payment_amount":   100000,
			"payment_due_date": time.Now().AddDate(0, 1, 0).Format(time.DateOnly),
			"fee_rate":         0.02,
			"tax_rate":         0.10,
		}
		invoiceBody, _ := json.Marshal(invoiceReq)

		req, _ := http.NewRequest(
			http.MethodPost,
			server.URL+"/api/invoices",
			bytes.NewBuffer(invoiceBody),
		)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		defer func() { _ = resp.Body.Close() }()
	})
}

func TestE2E_InvalidRequests(t *testing.T) {
	// テスト用DBのセットアップ
	db := setupTestDB(t)

	// テストデータの作成
	email, _ := setupTestData(t, db)

	// テスト用の設定
	cfg := &config.Config{
		JWTSecret: "test-secret-key-for-e2e",
	}

	// サーバーのセットアップ
	server := setupRouter(db, cfg)
	defer server.Close()

	// ログインしてトークンを取得
	loginReq := map[string]string{
		"email":    email,
		"password": "testpassword",
	}
	loginBody, _ := json.Marshal(loginReq)

	resp, _ := http.Post(
		server.URL+"/api/login",
		"application/json",
		bytes.NewBuffer(loginBody),
	)
	var loginResp map[string]string
	err := json.NewDecoder(resp.Body).Decode(&loginResp)
	assert.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	token := loginResp["token"]

	t.Run("E2E - 不正な請求書作成リクエスト", func(t *testing.T) {
		// 必須フィールド不足
		invoiceReq := map[string]interface{}{
			"payment_amount": 100000,
		}
		invoiceBody, _ := json.Marshal(invoiceReq)

		req, _ := http.NewRequest(
			http.MethodPost,
			server.URL+"/api/invoices",
			bytes.NewBuffer(invoiceBody),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		defer func() { _ = resp.Body.Close() }()
	})
}
