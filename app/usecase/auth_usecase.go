package usecase

import (
	"context"
	"errors"

	"github.com/ijufumi/practice-202512/app/config"
	"github.com/ijufumi/practice-202512/app/domain/repository"
	"github.com/ijufumi/practice-202512/app/util"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (string, error)
}

type authUsecase struct {
	userRepository repository.UserRepository
	config         *config.Config
}

func NewAuthUsecase(userRepository repository.UserRepository, cfg *config.Config) AuthUsecase {
	return &authUsecase{
		userRepository: userRepository,
		config:         cfg,
	}
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, error) {
	db, err := util.GetDB(ctx)
	if err != nil {
		return "", err
	}

	user, err := u.userRepository.FindByEmail(db, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid email or password")
		}

		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := util.GenerateJWT(user.ID, u.config.JWTSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
