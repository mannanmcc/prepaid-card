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

type RefundAmountCannotBeMoreThanCapturedAmountError struct {
	amount float64
}

func (a *RefundAmountCannotBeMoreThanCapturedAmountError) Error() string {
	return fmt.Sprintf("Oops, can not refund this amount")
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

	transaction.ReFund(amount)
	//topup the refunded amount
	a.Topup(amount)

	return nil
}
