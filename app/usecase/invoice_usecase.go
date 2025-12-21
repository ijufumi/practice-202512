package usecase

import (
	"context"
	"time"

	"github.com/ijufumi/practice-202512/app/domain/models"
	"github.com/ijufumi/practice-202512/app/domain/repository"
	"github.com/ijufumi/practice-202512/app/domain/value"
	"github.com/ijufumi/practice-202512/app/util"
)

type InvoiceUsecase interface {
	CreateInvoice(ctx context.Context, clientID string, issueDate time.Time, paymentAmount int, paymentDueDate time.Time) (*models.Invoice, error)
	GetInvoicesByPaymentDueDateRange(ctx context.Context, startDate, endDate *time.Time, offset, limit int) ([]*models.Invoice, error)
}

type invoiceUsecase struct {
	invoiceRepository repository.InvoiceRepository
	userRepository    repository.UserRepository
}

func NewInvoiceUsecase(invoiceRepository repository.InvoiceRepository, userRepository repository.UserRepository) InvoiceUsecase {
	return &invoiceUsecase{
		invoiceRepository: invoiceRepository,
		userRepository:    userRepository,
	}
}

func (u *invoiceUsecase) CreateInvoice(ctx context.Context, clientID string, issueDate time.Time, paymentAmount int, paymentDueDate time.Time) (*models.Invoice, error) {
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

	// 手数料率と消費税率の定義
	const feeRate = 0.04
	const taxRate = 0.10

	// 手数料計算: 支払金額 * 4%
	fee := int(float64(paymentAmount) * feeRate)

	// 消費税計算: 手数料 * 10%
	tax := int(float64(fee) * taxRate)

	// 請求金額計算: 支払金額 + 手数料 + 消費税
	invoiceAmount := paymentAmount + fee + tax

	invoice := &models.Invoice{
		CompanyID:      user.CompanyID,
		ClientID:       clientID,
		IssueDate:      issueDate,
		PaymentAmount:  paymentAmount,
		Fee:            fee,
		FeeRate:        feeRate,
		Tax:            tax,
		TaxRate:        taxRate,
		InvoiceAmount:  invoiceAmount,
		PaymentDueDate: paymentDueDate,
		Status:         value.InvoiceStatusUnprocessed,
	}

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
