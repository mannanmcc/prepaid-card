package models

import (
	"errors"
	"time"
)

//Transaction - represent the transaction which is captured from reserved
type Transaction struct {
	ID                   int
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

const STATUS_REFUNDED = "STATUS_REFUNDED"

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
		transaction.Status = STATUS_REFUNDED
	}

	transaction.Balance = transaction.Balance - amount
	return nil
}
