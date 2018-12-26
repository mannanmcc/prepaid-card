package handlers

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/mannanmcc/prepaid-card/models"
)

//CardHolderCommandInterface - list of required command for each type transaction
type CardHolderCommandInterface interface {
	TopupCommand(data []string, env Env) ([]string, error)
}

//CardHolderCommand - command type for card holder
type CardHolderCommand struct{}

//Toptup - toptup the card
func (command *CardHolderCommand) Toptup(cardNumber string, amount float64, db *gorm.DB) error {
	if amount < 0 {
		return errors.New("Negative number can not used for topup")
	}

	accountRepo := models.AccountRepository{Db: db}
	account, err := accountRepo.FindByCardNumber(cardNumber)
	if err != nil {
		return err
	}

	if account.Status != ACCOUNT_STATUS_ACTIVE {
		return errors.New("The account is inactive")
	}

	account.Topup(amount)
	if err := accountRepo.UpdateAccount(account); err != nil {
		return err
	}

	return nil
}
