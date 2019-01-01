package handlers

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mannanmcc/prepaid-card/models"
)

//Command - list of required command for each type transaction
type CapturedCommandInterface interface {
	captureFund(transactionReq *TransactionRequest, db *gorm.DB) error
	refund(refundReq *TransactionRequest, db *gorm.DB) error
}

type CaptureCommand struct{}

func (c *CaptureCommand) captureFund(transactionReq *TransactionRequest, db *gorm.DB) error {
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
	transaction := models.Transaction{
		CardNumber:           blockedTransaction.CardNumber,
		TransactionID:        transactionID,
		BlockedTransactionID: blockedTransaction.TransactionID,
		MerchantID:           blockedTransaction.MerchantID,
		Amount:               transactionReq.amount,
		Balance:              transactionReq.amount,
		Status:               models.StatusCaptured,
		CapturedAt:           time.Now(),
	}

	err = models.DoInTransaction(
		func(tx *gorm.DB) error {
			return tx.Create(&transaction).Error
		},
		func(tx *gorm.DB) error {
			return tx.Table("blocked_transactions").Where("id = ?", blockedTransaction.ID).Updates(map[string]interface{}{
				"balance": blockedTransaction.Balance,
			}).Error

		},
		db,
	)

	return err
}

func (c *CaptureCommand) refund(refundReq *TransactionRequest, db *gorm.DB) error {
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
	err = models.DoInTransaction(
		func(tx *gorm.DB) error {
			//use interface instead of model to get working with zero value update
			return tx.Table("transactions").Where("id = ?", transaction.ID).Updates(map[string]interface{}{
				"balance": transaction.Balance,
			}).Error
		},
		func(tx *gorm.DB) error {
			return tx.Table("account").Where("id = ?", account.ID).Update(map[string]interface{}{
				"balance": account.Balance,
			}).Error
		},
		db,
	)

	return err
}
