package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

type AccountRepositoryInterface interface {
	TopupAccount(account Account) (Account, error)
	FindByAccountNumber(accountNumber int) (*Account, error)
}

type AccountRepository struct {
	Db *gorm.DB
}

//FindByAccountNumber - returns account
func (repo *AccountRepository) FindByAccountNumber(accountNumber int) (Account, error) {
	var account Account
	res := repo.Db.Find(&account, &Account{AccountNumber: accountNumber})

	if res.RecordNotFound() {
		return account, errors.New(fmt.Sprintf("account not found with account number : %d", accountNumber))
	}

	return account, nil
}

//UpdateAccount - updates the account
func (repo *AccountRepository) UpdateAccount(account Account) (Account, error) {
	id := repo.Db.Save(&account)

	if id == nil {
		return account, errors.New("account saving failed")
	}

	return account, nil
}
