package usecase

import (
	"context"
	"time"

	"github.com/ijufumi/practice-202512/app/config"
	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/domain/repository"
	"github.com/ijufumi/practice-202512/app/domain/value"
	"github.com/ijufumi/practice-202512/app/util"
	"github.com/shopspring/decimal"
)

type InvoiceUsecase interface {
	CreateInvoice(ctx context.Context, clientID string, issueDate time.Time, paymentAmount decimal.Decimal, paymentDueDate time.Time) (*models.Invoice, error)
	GetInvoicesByPaymentDueDateRange(ctx context.Context, startDate, endDate *time.Time, offset, limit int) ([]*models.Invoice, error)
}

type invoiceUsecase struct {
	invoiceRepository repository.InvoiceRepository
	userRepository    repository.UserRepository
	config            *config.Config
}

func NewInvoiceUsecase(invoiceRepository repository.InvoiceRepository, userRepository repository.UserRepository) InvoiceUsecase {
	return &invoiceUsecase{
		invoiceRepository: invoiceRepository,
		userRepository:    userRepository,
		config:            config.Load(),
	}
}

func (u *invoiceUsecase) CreateInvoice(ctx context.Context, clientID string, issueDate time.Time, paymentAmount decimal.Decimal, paymentDueDate time.Time) (*models.Invoice, error) {
	db, err := util.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := util.GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepository.FindByID(db, userID)
	if err != nil {
		return nil, err
	}

	invoice := &models.Invoice{
		CompanyID:      user.CompanyID,
		ClientID:       clientID,
		IssueDate:      issueDate,
		PaymentAmount:  paymentAmount,
		PaymentDueDate: paymentDueDate,
		Status:         value.InvoiceStatusUnprocessed,
	}

	// domain/models の計算メソッドを使用
	invoice.CalculateFee(u.config.FeeRate)
	invoice.CalculateTax(u.config.TaxRate)
	invoice.CalculateInvoiceAmount()

	if err := u.invoiceRepository.Create(db, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (u *invoiceUsecase) GetInvoicesByPaymentDueDateRange(ctx context.Context, startDate, endDate *time.Time, offset, limit int) ([]*models.Invoice, error) {
	db, err := util.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	// デフォルト値の設定
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 100
	}

	return u.invoiceRepository.FindByPaymentDueDateRange(db, startDate, endDate, offset, limit)
}
