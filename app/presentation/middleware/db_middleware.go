package middleware

import (
	"github.com/ijufumi/practice-202512/app/util"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

// DBMiddleware adds gorm.DB instance to request context
func DBMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := util.SetDB(req.Context(), db.WithContext(req.Context()))
			c.SetRequest(req.WithContext(ctx))
			return next(c)
		}
	}
}
