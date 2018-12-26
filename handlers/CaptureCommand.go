package handlers

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mannanmcc/prepaid-card/models"
)

//Command - list of required command for each type transaction
type CapturedCommandInterface interface {
	CaptureFundCommand(data []string, env Env) ([]string, error)
	RefundCommand(data []string, env Env) ([]string, error)
}

type CaptureCommand struct{}

//CaptureFund - capture blocked
func (c *CaptureCommand) CaptureFund(transactionReq *TransactionRequest, db *gorm.DB) error {
	blockTransactionRepo := models.BlockedTransactionRepository{Db: db}
	blockedTransaction, err := blockTransactionRepo.FindByTransactionID(transactionReq.transactionId)
	if err != nil {
		return err
	}

	//capture the fund which was blocked to be captured and remove balance to make sure it can not be captured again for the same amount
	if err := blockedTransaction.CaptureFund(transactionReq.amount); err != nil {
		return err
	}

	transactionID := models.GenerateTransactionId(RANDOM_KEY_LENGTH)
	transactionRepo := models.TransactionRepository{Db: db}
	//create a transaction based on blocked transaction
	transaction := models.Transaction{
		CardNumber:           blockedTransaction.CardNumber,
		TransactionID:        transactionID,
		BlockedTransactionID: blockedTransaction.TransactionID,
		MerchantID:           blockedTransaction.MerchantID,
		Amount:               transactionReq.amount,
		Balance:              transactionReq.amount,
		Status:               models.STATUS_CAPTURED,
		CapturedAt:           time.Now(),
	}

	if _, err := transactionRepo.Create(transaction); err != nil {
		return err
	}

	if err := blockTransactionRepo.Update(blockedTransaction); err != nil {
		return err
	}

	return nil
}

//Refund - refund the captured fund
func (c *CaptureCommand) Refund(refundReq *TransactionRequest, db *gorm.DB) error {
	transactionRepo := models.TransactionRepository{Db: db}
	transaction, err := transactionRepo.FindByTransactionID(refundReq.transactionId)
	if err != nil {
		return err
	}

	accountRepo := models.AccountRepository{Db: db}
	account, err := accountRepo.FindByCardNumber(transaction.CardNumber)
	if err != nil {
		return err
	}

	if err = transaction.Refund(refundReq.amount); err != nil {
		return err
	}

	account.Topup(refundReq.amount)

	//update the account
	if err := accountRepo.Update(account); err != nil {
		return err
	}

	err = transactionRepo.Update(transaction)
	return err
}
