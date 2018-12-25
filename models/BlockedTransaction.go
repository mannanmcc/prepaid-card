package models

import (
	"errors"
	"fmt"
	"time"
)

const STATUS_BLOCKED = "BLOCKED"
const STATUS_CAPTURED = "CAPTURED"
const STATUS_REVERSED = "REVERSED"

//todo - add following later on startDate
type BlockedTransaction struct {
	ID            int
	AccountNumber int    `gorm:"column:account_number"`
	TransactionID string `gorm:"column:transaction_id"`
	Amount        float64
	MerchantID    string
	Reason        string
	BlockedAt     time.Time `json:"BlockedAt" gorm:"column:blocked_at;type:datetime;default:CURRENT_TIMESTAMP"`
	Status        string    `gorm:"column:status"`
}

func (bt *BlockedTransaction) TableName() string {
	return "blocked_transactions"
}

type CannotBlockMoreThanCurrentBalance struct {
	Amount float64
}

func (a *CannotBlockMoreThanCurrentBalance) Error() string {
	return fmt.Sprintf("Cannot block amount %f as currently account does not have sufficient money", a.Amount)
}

type CannotBlockAmount struct {
	Amount float64
}

func (a *CannotBlockAmount) Error() string {
	return fmt.Sprintf("Cannot block amount: %f", a.Amount)
}

//Capture - check capture with blocked amount and decrease the blocked if success
func (bt *BlockedTransaction) CaptureFund(amount float64) error {
	if bt.Amount < amount {
		return errors.New("Cannot capture amount which is more than remaining blocked amount")
	}

	bt.Amount = bt.Amount - amount
	//Changed status to captured if capturing full amount
	if bt.Amount == amount {
		bt.Status = STATUS_CAPTURED
	}

	return nil
}

//Reverse - check capture with blocked amount and decrease the blocked if success
func (bt *BlockedTransaction) Reverse(amount float64) error {
	if bt.Amount < amount {
		return errors.New("Cannot capture amount which is more than captured amount")
	}

	fmt.Printf("Cannot capture amount which is more than captured amount 1")
	bt.Amount = bt.Amount - amount
	//Changed status to captured if capturing full amount
	if bt.Amount == amount {
		fmt.Printf("Cannot capture amount which is more than captured amount 2")
		bt.Status = STATUS_REVERSED
	}

	return nil
}
