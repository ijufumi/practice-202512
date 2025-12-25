package main

import (
	"log"
	"time"

	"github.com/ijufumi/practice-202512/app/config"
	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/domain/value"
	"github.com/ijufumi/practice-202512/app/infrastructure/database"
	"github.com/ijufumi/practice-202512/app/infrastructure/database/gateway"
	"github.com/shopspring/decimal"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// 設定の読み込み
	cfg := config.Load()

	// データベース接続
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepository := gateway.NewUserRepository()
	companyRepository := gateway.NewCompanyRepository()
	clientRepository := gateway.NewClientRepository()
	clientBankAccountRepository := gateway.NewClientBankAccountRepository()
	invoiceRepository := gateway.NewInvoiceRepository()
	err = db.Transaction(func(tx *gorm.DB) error {
		company := &models.Company{
			CorporateName:      "test corporation",
			RepresentativeName: "test representative",
			PhoneNumber:        "000-0000-0000",
			PostalCode:         "000-0000",
			Address:            "test address",
		}
		if err := companyRepository.Create(tx, company); err != nil {
			return err
		}
		passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user := &models.User{
			CompanyID: company.ID,
			Name:      "admin",
			Email:     "admin@localhost.ai",
			Password:  string(passwordHash),
		}
		if err := userRepository.Create(tx, user); err != nil {
			return err
		}

		client := &models.Client{
			CompanyID:          company.ID,
			CorporateName:      "test corporation",
			RepresentativeName: "test representative",
			PhoneNumber:        "000-0000-0000",
			PostalCode:         "000-0000",
			Address:            "test address",
		}
		if err := clientRepository.Create(tx, client); err != nil {
			return err
		}
		clientBankAccount := &models.ClientBankAccount{
			ClientID:      client.ID,
			BankName:      "test bank",
			BranchName:    "test branch",
			AccountNumber: "0000000000000",
			AccountName:   "test account",
		}
		if err := clientBankAccountRepository.Create(tx, clientBankAccount); err != nil {
			return err
		}
		invoice := &models.Invoice{
			CompanyID:      company.ID,
			ClientID:       client.ID,
			IssueDate:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
			PaymentAmount:  decimal.NewFromInt(10000),
			PaymentDueDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.Local),
			Fee:            decimal.NewFromInt(1000),
			FeeRate:        decimal.NewFromFloat(0.01),
			Tax:            decimal.NewFromInt(1000),
			TaxRate:        decimal.NewFromFloat(0.01),
			InvoiceAmount:  decimal.NewFromInt(10000),
			Status:         value.InvoiceStatusProcessed,
		}
		if err := invoiceRepository.Create(tx, invoice); err != nil {
			return err
		}
		return nil
	})
}
