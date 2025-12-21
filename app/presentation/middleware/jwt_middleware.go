package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/ijufumi/practice-202512/app/config"
	"github.com/ijufumi/practice-202512/app/util"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			log.Println(authHeader)
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
				})
			}

			tokenString := parts[1]
			claims, err := util.ValidateJWT(tokenString, cfg.JWTSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
			}
			ctx := c.Request().Context()
			ctx = util.SetUserID(ctx, claims.UserID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
