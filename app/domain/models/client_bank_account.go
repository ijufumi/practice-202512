package models

import (
	"github.com/ijufumi/practice-202512/app/infrastructure/database/dao"

	"time"
)

type ClientBankAccount struct {
	ID            string
	ClientID      string
	BankName      string
	BranchName    string
	AccountNumber string
	AccountName   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (c *ClientBankAccount) ToDAO() *dao.ClientBankAccount {
	return &dao.ClientBankAccount{
		ID:            c.ID,
		ClientID:      c.ClientID,
		BankName:      c.BankName,
		BranchName:    c.BranchName,
		AccountNumber: c.AccountNumber,
		AccountName:   c.AccountName,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

func ClientBankAccountFromDAO(daoAccount *dao.ClientBankAccount) *ClientBankAccount {
	return &ClientBankAccount{
		ID:            daoAccount.ID,
		ClientID:      daoAccount.ClientID,
		BankName:      daoAccount.BankName,
		BranchName:    daoAccount.BranchName,
		AccountNumber: daoAccount.AccountNumber,
		AccountName:   daoAccount.AccountName,
		CreatedAt:     daoAccount.CreatedAt,
		UpdatedAt:     daoAccount.UpdatedAt,
	}
}
