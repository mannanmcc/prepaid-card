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
func (command *CardHolderCommand) Toptup(topupRequest *TopupRequest, db *gorm.DB) error {
	if topupRequest.amount < 0 {
		return errors.New("Negative number can not used for topup")
	}

	accountRepo := models.AccountRepository{Db: db}
	account, err := accountRepo.FindByCardNumber(topupRequest.cardNumber)
	if err != nil {
		return err
	}

	if account.Status != AccountStatusInActive {
		return errors.New("The account is inactive")
	}

	account.Topup(topupRequest.amount)
	if err := accountRepo.Update(account); err != nil {
		return err
	}

	return nil
}
