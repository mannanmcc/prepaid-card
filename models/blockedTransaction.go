package models

import (
	"errors"
	"fmt"
	"time"
)

const StatusBlocked = "BLOCKED"
const StatusCaptured = "CAPTURED"
const StatusReversed = "REVERSED"

//todo - add following later on startDate
type BlockedTransaction struct {
	ID                  int    `gorm:"primary_key"`
	CardNumber          string `gorm:"column:card_number"`
	TransactionID       string `gorm:"column:transaction_id"`
	ParentTransactionID string `gorm:"column:parent_transaction_id"`
	Balance             float64
	Amount              float64
	MerchantID          string
	Reason              string
	BlockedAt           time.Time
	UpdatedAt           time.Time
	Status              string `gorm:"column:status"`
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

//CaptureFund - check capture with blocked amount and decrease the blocked if success
func (bt *BlockedTransaction) CaptureFund(amount float64) error {
	if bt.Balance < amount {
		return errors.New("Cannot capture amount which is more than remaining blocked amount")
	}

	bt.Balance = bt.Balance - amount
	return nil
}

//Reverse - check capture with blocked amount and decrease the blocked if success
func (bt *BlockedTransaction) Reverse(amount float64) error {
	if bt.Balance < amount {
		return errors.New("Cannot capture amount which is more than captured amount")
	}

	//Changed status to captured if capturing full amount
	if bt.Balance == amount {
		bt.Status = StatusReversed
	}

	bt.Balance = bt.Balance - amount

	return nil
}
