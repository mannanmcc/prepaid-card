package models

import (
	"errors"
	"time"
)

//Transaction - represent the transaction which is captured from reserved
type Transaction struct {
	ID                   int `gorm:"primary_key"`
	CardNumber           string
	BlockedTransactionID string
	TransactionID        string
	ParentTransactionID  string
	MerchantID           string
	Amount               float64
	Balance              float64
	Status               string
	CapturedAt           time.Time
}

const statusRefund = "REFUND"

//TableName - db table name
func (transaction *Transaction) TableName() string {
	return "transactions"
}

//Refund - check capture with blocked amount and decrease the blocked if success
func (transaction *Transaction) Refund(amount float64) error {
	if transaction.Balance < amount {
		return errors.New("Cannot refund amount which is more than captured amount")
	}

	if transaction.Balance == amount {
		transaction.Status = statusRefund
	}

	transaction.Balance = transaction.Balance - amount
	return nil
}
