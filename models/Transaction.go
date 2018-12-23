package models

import (
	"errors"
	"fmt"
	"time"
)

//todo - add following later on startDate
type Transaction struct {
	ID                   int
	AccountNumber        int
	BlockedTransactionID string
	MerchantID           string
	Amount               float64
	Status               string
	TransactionID        string
	CapturedAt           time.Time
}

const STATUS_REFUNDED = "STATUS_REFUNDED"

func (bt *Transaction) TableName() string {
	return "transactions"
}

//Capture - check capture with blocked amount and decrease the blocked if success
func (transaction *Transaction) ReFund(amount float64) error {
	if transaction.Amount < amount {
		return errors.New("Cannot capture amount which is more than captured amount")
	}

	fmt.Printf("Cannot capture amount which is more than captured amount 1")
	transaction.Amount = transaction.Amount - amount
	//Changed status to captured if capturing full amount
	if transaction.Amount == amount {
		fmt.Printf("Cannot capture amount which is more than captured amount 2")
		transaction.Status = STATUS_REFUNDED
	}

	return nil
}
