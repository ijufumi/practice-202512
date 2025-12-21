package handler

import (
	"net/http"

	"github.com/ijufumi/practice-202512/app/presentation/models"
	"github.com/ijufumi/practice-202512/app/usecase"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request body"))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
	}

	token, err := h.authUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
	})
}
