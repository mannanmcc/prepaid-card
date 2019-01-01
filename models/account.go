package models

import (
	"errors"
	"fmt"
	"time"
)

const (
	AccountStatusInActive = "INACTIVE"
	AccountStatusActive   = "ACTIVE"
	AccountStatusDisabled = "DISABLED"
)

//todo - add following later on startDate
type Account struct {
	ID                int `gorm:"primary_key"`
	Balance           float64
	AccountHolderName string
	CardNumber        string
	SortCode          int
	Status            string
	Address           string
	Postcode          string
	DateOfBirth       time.Time
}

func (account *Account) TableName() string {
	return "account"
}

type AccountNotExists struct {
	CardNumber int
}

func (a *AccountNotExists) Error() string {
	return fmt.Sprintf("account number '%d' not exists", a.CardNumber)
}

type NotEnoughMoneyToPayError struct {
	amount float64
}

func (a *NotEnoughMoneyToPayError) Error() string {
	return fmt.Sprintf("account number does not have enough money to pay")
}

type RefundAmountCannotBeMoreThanCapturedAmountError struct {
	amount float64
}

func (a *RefundAmountCannotBeMoreThanCapturedAmountError) Error() string {
	return fmt.Sprintf("Oops, can not refund this amount")
}

//ReverseAmountCannotBeMoreThanCapturedAmountError - reverse the authorise amount
type ReverseAmountCannotBeMoreThanCapturedAmountError struct {
	amount float64
}

func (a *ReverseAmountCannotBeMoreThanCapturedAmountError) Error() string {
	return fmt.Sprintf("Oops, can not reverse this amount")
}

//Topup - add amount with exisiting amount in the account
func (a *Account) Topup(amount float64) {
	a.Balance = a.Balance + amount
}

//AuthoriseAmount - return the transaction amount against account balance
func (a *Account) AuthoriseAmount(amount float64) (bool, error) {
	if amount > a.Balance {
		return false, &CannotBlockMoreThanCurrentBalance{}
	}

	a.Balance = a.Balance - amount
	return true, nil
}

//Refund - return the transaction amount against account balance
func (a *Account) Refund(transaction *Transaction, amount float64) error {
	if amount > transaction.Amount {
		return &RefundAmountCannotBeMoreThanCapturedAmountError{}
	}

	transaction.Refund(amount)
	//topup the refunded amount
	a.Topup(amount)

	return nil
}

//Reverse - return the transaction amount against account balance
func (a *Account) ReverseAuthorisedAmount(transaction *BlockedTransaction, amount float64) error {
	if err := transaction.Reverse(amount); err != nil {
		return errors.New("oops, transaction can not be reverse")
	}

	//topup the refunded amount
	a.Topup(amount)

	return nil
}
