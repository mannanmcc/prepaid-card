package handlers

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mannanmcc/prepaid-card/models"
)

//Command - list of required command for each type transaction
type CommandInterface interface {
	authorisationCommand(authReq authorisationRequestBody, db *gorm.DB) error
	reverseCommand(transactionReq *TransactionRequest, db *gorm.DB) error
}

type Command struct{}

func (c *Command) authorisationCommand(authReq authorisationRequestBody, db *gorm.DB) error {
	accountRepo := models.AccountRepository{Db: db}
	account, err := accountRepo.FindByCardNumber(authReq.cardNumber)
	if err != nil {
		return err
	}

	_, err = account.AuthoriseAmount(authReq.amount)
	if err != nil {
		return err
	}

	transactionID := models.GenerateTransactionId(RANDOM_KEY_LENGTH)
	blockTransaction := models.BlockedTransaction{
		TransactionID: transactionID,
		CardNumber:    account.CardNumber,
		Amount:        authReq.amount,
		Balance:       authReq.amount,
		MerchantID:    authReq.merchantId,
		Reason:        authReq.reason,
		BlockedAt:     time.Now(),
		Status:        models.StatusBlocked,
	}

	err = models.DoInTransaction(
		func(tx *gorm.DB) error {
			return tx.Create(&blockTransaction).Error
		},
		func(tx *gorm.DB) error {
			return tx.Table("account").Where("id = ?", account.ID).Updates(map[string]interface{}{
				"balance": account.Balance,
			}).Error
		},
		db,
	)

	return err
}

func (c *Command) reverseCommand(transactionReq *TransactionRequest, db *gorm.DB) error {
	transactionRepo := models.BlockedTransactionRepository{Db: db}
	transaction, err := transactionRepo.FindByTransactionID(transactionReq.transactionId)
	if err != nil {
		return err
	}

	accountRepo := models.AccountRepository{Db: db}
	account, err := accountRepo.FindByCardNumber(transaction.CardNumber)
	if err != nil {
		return err
	}

	if err = transaction.Reverse(transactionReq.amount); err != nil {
		return err
	}

	account.Topup(transactionReq.amount)

	err = models.DoInTransaction(
		func(tx *gorm.DB) error {
			//use interface instead of model to get working with zero value update
			return tx.Table("blocked_transactions").Where("transaction_id = ?", transactionReq.transactionId).Updates(map[string]interface{}{
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
