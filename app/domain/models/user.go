package models

import (
	"github.com/ijufumi/practice-202512/app/infrastructure/database/dao"

	"time"
)

type User struct {
	ID        string
	CompanyID string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) ToDAO() *dao.User {
	return &dao.User{
		ID:        u.ID,
		CompanyID: u.CompanyID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func UserFromDAO(daoUser *dao.User) *User {
	return &User{
		ID:        daoUser.ID,
		CompanyID: daoUser.CompanyID,
		Name:      daoUser.Name,
		Email:     daoUser.Email,
		Password:  daoUser.Password,
		CreatedAt: daoUser.CreatedAt,
		UpdatedAt: daoUser.UpdatedAt,
	}
}
