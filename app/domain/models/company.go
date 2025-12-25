package models

import (
	"github.com/ijufumi/practice-202512/app/infrastructure/database/entities"

	"time"
)

type Company struct {
	ID                 string
	CorporateName      string
	RepresentativeName string
	PhoneNumber        string
	PostalCode         string
	Address            string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (c *Company) ToDAO() *entities.Company {
	return &entities.Company{
		ID:                 c.ID,
		CorporateName:      c.CorporateName,
		RepresentativeName: c.RepresentativeName,
		PhoneNumber:        c.PhoneNumber,
		PostalCode:         c.PostalCode,
		Address:            c.Address,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}

func CompanyFromDAO(daoCompany *entities.Company) *Company {
	return &Company{
		ID:                 daoCompany.ID,
		CorporateName:      daoCompany.CorporateName,
		RepresentativeName: daoCompany.RepresentativeName,
		PhoneNumber:        daoCompany.PhoneNumber,
		PostalCode:         daoCompany.PostalCode,
		Address:            daoCompany.Address,
		CreatedAt:          daoCompany.CreatedAt,
		UpdatedAt:          daoCompany.UpdatedAt,
	}
}
