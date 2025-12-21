package handler

import (
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

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")
	offsetStr := c.QueryParam("offset")
	limitStr := c.QueryParam("limit")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		_startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid start_date format. Use YYYY-MM-DD"))
		}
		startDate = &_startDate
	}

	if endDateStr != "" {
		_endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid end_date format. Use YYYY-MM-DD"))
		}
		endDate = &_endDate
	}

	// offset と limit のパース（デフォルト値: offset=0, limit=100）
	offset := 0
	limit := 100

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid offset parameter"))
		}
		offset = parsedOffset
	}

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid limit parameter"))
		}
		limit = parsedLimit
	}

	invoices, err := h.invoiceUsecase.GetInvoicesByPaymentDueDateRange(ctx, startDate, endDate, offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get invoices"))
	}

	responses := models.FromInvoiceDomainModels(invoices)

	return c.JSON(http.StatusOK, responses)
}
