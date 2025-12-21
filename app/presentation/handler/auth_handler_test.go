package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/ijufumi/practice-202512/app/presentation/models"
	usecase "github.com/ijufumi/practice-202512/app/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	return e
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("ログイン成功", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockAuthUsecase(t)

		mockUsecase.EXPECT().Login(mock.Anything, "test@example.com", "password123").
			Return("test-jwt-token", nil)

		handler := NewAuthHandler(mockUsecase)

		reqBody := `{"email":"test@example.com","password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.LoginResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test-jwt-token", response.Token)
	})

	t.Run("不正なリクエストボディ", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockAuthUsecase(t)

		handler := NewAuthHandler(mockUsecase)

		reqBody := `{"invalid":"json"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("バリデーションエラー - Emailなし", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockAuthUsecase(t)

		handler := NewAuthHandler(mockUsecase)

		reqBody := `{"password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("バリデーションエラー - 不正なEmail形式", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockAuthUsecase(t)

		handler := NewAuthHandler(mockUsecase)

		reqBody := `{"email":"invalid-email","password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("バリデーションエラー - パスワードなし", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockAuthUsecase(t)

		handler := NewAuthHandler(mockUsecase)

		reqBody := `{"email":"test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("認証エラー", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockAuthUsecase(t)

		mockUsecase.EXPECT().Login(mock.Anything, "test@example.com", "wrongpassword").
			Return("", errors.New("invalid email or password"))

		handler := NewAuthHandler(mockUsecase)

		reqBody := `{"email":"test@example.com","password":"wrongpassword"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid email or password", response["error"])
	})

	t.Run("内部サーバーエラー", func(t *testing.T) {
		e := setupEcho()
		mockUsecase := usecase.NewMockAuthUsecase(t)

		mockUsecase.EXPECT().Login(mock.Anything, "test@example.com", "password123").
			Return("", errors.New("internal error"))

		handler := NewAuthHandler(mockUsecase)

		reqBody := `{"email":"test@example.com","password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
