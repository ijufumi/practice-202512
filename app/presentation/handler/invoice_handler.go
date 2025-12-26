package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ijufumi/practice-202512/app/presentation/models"
	"github.com/ijufumi/practice-202512/app/usecase"

	"github.com/labstack/echo/v4"
)

type InvoiceHandler struct {
	invoiceUsecase usecase.InvoiceUsecase
}

func NewInvoiceHandler(invoiceUsecase usecase.InvoiceUsecase) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceUsecase: invoiceUsecase,
	}
}

func (h *InvoiceHandler) CreateInvoice(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.CreateInvoiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request body"))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
	}

	// 日付のパース
	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid issue_date format. Use YYYY-MM-DD"))
	}

	paymentDueDate, err := time.Parse("2006-01-02", req.PaymentDueDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid payment_due_date format. Use YYYY-MM-DD"))
	}

	invoice, err := h.invoiceUsecase.CreateInvoice(ctx, req.ClientID, issueDate, req.PaymentAmount, paymentDueDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to create invoice"))
	}

	response := models.FromInvoiceDomainModel(invoice)

	return c.JSON(http.StatusCreated, response)
}

func (h *InvoiceHandler) GetInvoices(c echo.Context) error {
	ctx := c.Request().Context()

	// クエリパラメータのパース
	startDate, err := parseOptionalDate(c.QueryParam("start_date"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid start_date format. Use YYYY-MM-DD"))
	}

	endDate, err := parseOptionalDate(c.QueryParam("end_date"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid end_date format. Use YYYY-MM-DD"))
	}

	offset, err := parseOffset(c.QueryParam("offset"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid offset parameter"))
	}

	limit, err := parseLimit(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid limit parameter"))
	}

	invoices, err := h.invoiceUsecase.GetInvoicesByPaymentDueDateRange(ctx, startDate, endDate, offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get invoices"))
	}

	responses := models.FromInvoiceDomainModels(invoices)

	return c.JSON(http.StatusOK, responses)
}

// parseOptionalDate はオプショナルな日付文字列をパースします
func parseOptionalDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	return &parsedDate, nil
}

// parseOffset はオフセット文字列をパースします（デフォルト: 0）
func parseOffset(offsetStr string) (int, error) {
	if offsetStr == "" {
		return DefaultOffset, nil
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return 0, err
	}
	if offset <= 0 {
		return 0, errors.New("invalid limit parameter")
	}

	return offset, nil
}

// parseLimit はリミット文字列をパースします（デフォルト: 100）
func parseLimit(limitStr string) (int, error) {
	if limitStr == "" {
		return DefaultLimit, nil
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, err
	}
	if limit <= 0 {
		return 0, errors.New("invalid limit parameter")
	}

	return limit, nil
}
