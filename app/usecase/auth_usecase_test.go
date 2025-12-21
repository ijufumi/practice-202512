package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/ijufumi/practice-202512/app/config"
	"github.com/ijufumi/practice-202512/app/domain/models"
	repository "github.com/ijufumi/practice-202512/app/domain/repository/mocks"
	"github.com/ijufumi/practice-202512/app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupContext(t *testing.T) (context.Context, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	ctx := util.SetDB(context.Background(), db)
	return ctx, db
}

func TestAuthUsecase_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	t.Run("ログイン成功", func(t *testing.T) {
		ctx, _ := setupContext(t)
		mockRepo := repository.NewMockUserRepository(t)
		cfg := &config.Config{
			JWTSecret: "test-secret",
		}

		expectedUser := &models.User{
			ID:       "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
			Email:    "test@example.com",
			Password: string(hashedPassword),
		}

		mockRepo.EXPECT().FindByEmail(mock.Anything, "test@example.com").
			Return(expectedUser, nil)

		usecase := NewAuthUsecase(mockRepo, cfg)
		token, err := usecase.Login(ctx, "test@example.com", "password123")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("ユーザーが見つからない", func(t *testing.T) {
		ctx, _ := setupContext(t)
		mockRepo := repository.NewMockUserRepository(t)
		cfg := &config.Config{
			JWTSecret: "test-secret",
		}

		mockRepo.EXPECT().FindByEmail(mock.Anything, "notfound@example.com").
			Return(nil, gorm.ErrRecordNotFound)

		usecase := NewAuthUsecase(mockRepo, cfg)
		token, err := usecase.Login(ctx, "notfound@example.com", "password123")

		assert.Error(t, err)
		assert.Equal(t, "invalid email or password", err.Error())
		assert.Empty(t, token)
	})

	t.Run("パスワードが間違っている", func(t *testing.T) {
		ctx, _ := setupContext(t)
		mockRepo := repository.NewMockUserRepository(t)
		cfg := &config.Config{
			JWTSecret: "test-secret",
		}

		expectedUser := &models.User{
			ID:       "01HQZXFG0PJ9K8QXW7YM1N2ZXC",
			Email:    "test@example.com",
			Password: string(hashedPassword),
		}

		mockRepo.EXPECT().FindByEmail(mock.Anything, "test@example.com").
			Return(expectedUser, nil)

		usecase := NewAuthUsecase(mockRepo, cfg)
		token, err := usecase.Login(ctx, "test@example.com", "wrongpassword")

		assert.Error(t, err)
		assert.Equal(t, "invalid email or password", err.Error())
		assert.Empty(t, token)
	})

	t.Run("リポジトリエラー", func(t *testing.T) {
		ctx, _ := setupContext(t)
		mockRepo := repository.NewMockUserRepository(t)
		cfg := &config.Config{
			JWTSecret: "test-secret",
		}

		mockRepo.EXPECT().FindByEmail(mock.Anything, "test@example.com").
			Return(nil, errors.New("database error"))

		usecase := NewAuthUsecase(mockRepo, cfg)
		token, err := usecase.Login(ctx, "test@example.com", "password123")

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Empty(t, token)
	})

	t.Run("コンテキストにDBがない", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := repository.NewMockUserRepository(t)
		cfg := &config.Config{
			JWTSecret: "test-secret",
		}

		usecase := NewAuthUsecase(mockRepo, cfg)
		token, err := usecase.Login(ctx, "test@example.com", "password123")

		assert.Error(t, err)
		assert.Equal(t, "database connection not found in context", err.Error())
		assert.Empty(t, token)
	})
}
