package models

import (
	"github.com/ijufumi/practice-202512/app/infrastructure/database/entities"

	"time"
)

type Client struct {
	ID                 string
	CompanyID          string
	CorporateName      string
	RepresentativeName string
	PhoneNumber        string
	PostalCode         string
	Address            string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (c *Client) ToDAO() *entities.Client {
	return &entities.Client{
		ID:                 c.ID,
		CompanyID:          c.CompanyID,
		CorporateName:      c.CorporateName,
		RepresentativeName: c.RepresentativeName,
		PhoneNumber:        c.PhoneNumber,
		PostalCode:         c.PostalCode,
		Address:            c.Address,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}

func ClientFromDAO(daoClient *entities.Client) *Client {
	return &Client{
		ID:                 daoClient.ID,
		CompanyID:          daoClient.CompanyID,
		CorporateName:      daoClient.CorporateName,
		RepresentativeName: daoClient.RepresentativeName,
		PhoneNumber:        daoClient.PhoneNumber,
		PostalCode:         daoClient.PostalCode,
		Address:            daoClient.Address,
		CreatedAt:          daoClient.CreatedAt,
		UpdatedAt:          daoClient.UpdatedAt,
	}
}
