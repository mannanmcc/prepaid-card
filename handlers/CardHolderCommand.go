package handlers

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/mannanmcc/prepaid-card/models"
)

const (
	cardNumberMin = 10000000
	cardNumberMax = 99999999
	sortCode      = 200300
)

//CardHolderCommandInterface - list of required command for each type transaction
type CardHolderCommandInterface interface {
	TopupCommand(TopupRequest, *gorm.DB) error
	CreditCard(CreatCardRequest *gorm.DB) error
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

func generateUniqueRandomNumber(min, max int, db *gorm.DB) string {
	accountRepo := models.AccountRepository{Db: db}
	cardNumberFound := true
	var generatedCardNumber string

	for cardNumberFound == true {
		generatedCardNumber = strconv.Itoa(min + rand.Intn(max-min))
		_, err := accountRepo.FindByCardNumber(generatedCardNumber)
		if err != nil {
			cardNumberFound = false
		}
	}

	return generatedCardNumber
}

//CreateCard - create a new card
func (command *CardHolderCommand) CreateCard(createCardRequest *CreateCardRequest, db *gorm.DB) (models.Account, error) {
	newCardNumber := generateUniqueRandomNumber(cardNumberMin, cardNumberMax, db)
	account := models.Account{
		CardNumber:        newCardNumber,
		AccountHolderName: createCardRequest.fullName,
		SortCode:          sortCode,
		Address:           createCardRequest.address,
		Postcode:          createCardRequest.postcode,
		DateOfBirth:       createCardRequest.dateOfBirth,
		Status:            models.AccountStatusActive,
	}

	accountRepo := models.AccountRepository{Db: db}
	accountRepo.Create(&account)
	if account.ID == 0 {
		return account, errors.New("Sorry, account creation failed, please try again")
	}

	return account, nil
}
