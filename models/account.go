package models

import "fmt"

//todo - add following later on startDate
type Account struct {
	ID                int
	Balance           float64
	AccountHolderName string
	AccountNumber     int
	SortCode          string
	Status            string
}

func (account *Account) TableName() string {
	return "account"
}

type AccountNotExists struct {
	AccountNumber int
}

func (a *AccountNotExists) Error() string {
	return fmt.Sprintf("account number '%d' not exists", a.AccountNumber)
}

type NotEnoughMoneyToPayError struct {
	amount float64
}

func (a *NotEnoughMoneyToPayError) Error() string {
	return fmt.Sprintf("account number does not have enough money to pay")
}

//Topup - add amount with exisiting amount in the account
func (a *Account) Topup(amount float64) {
	a.Balance = a.Balance + amount
}

//AuthoriseAmount - return the transaction amount against account balance
func (a *Account) AuthoriseAmount(amount float64) (bool, error) {
	if amount > a.Balance {
		return false, &NotEnoughMoneyToPayError{}
	}

	//add a transaction against this account
	return true, nil
}
