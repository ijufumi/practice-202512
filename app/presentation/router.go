package presentation

import (
	"github.com/ijufumi/practice-202512/app/config"
	"github.com/ijufumi/practice-202512/app/presentation/handler"
	custommiddleware "github.com/ijufumi/practice-202512/app/presentation/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, cfg *config.Config, invoiceHandler *handler.InvoiceHandler, authHandler *handler.AuthHandler) *echo.Echo {
	e := echo.New()

	// バリデーション
	e.Validator = custommiddleware.NewCustomValidator()

	// ミドルウェア
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(custommiddleware.DBMiddleware(db))

	// 認証API
	api := e.Group("/api")
	api.POST("/login", authHandler.Login)

	// 請求書API（JWT認証が必要）
	invoices := api.Group("/invoices")
	invoices.Use(custommiddleware.JWTMiddleware(cfg))
	invoices.POST("", invoiceHandler.CreateInvoice)
	invoices.GET("", invoiceHandler.GetInvoices)

	return e
}
