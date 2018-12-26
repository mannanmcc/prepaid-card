package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

//AccountRepositoryInterface - repository interface for account
type AccountRepositoryInterface interface {
	FindByCardNumber(cardNumber string) (*Account, error)
	Update(account *Account) error
}

//AccountRepository - type for account repository
type AccountRepository struct {
	Db *gorm.DB
}

//FindByCardNumber - returns account
func (repo *AccountRepository) FindByCardNumber(cardNumber string) (*Account, error) {
	var account Account
	res := repo.Db.Find(&account, &Account{CardNumber: cardNumber})

	if res.RecordNotFound() {
		return &account, fmt.Errorf("Card details not found with card number : %s", cardNumber)
	}

	return &account, nil
}

//Update - updates the account
func (repo *AccountRepository) Update(account *Account) error {
	id := repo.Db.Save(&account)

	if id == nil {
		return errors.New("account saving failed")
	}

	return nil
}
