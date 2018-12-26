package handlers

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mannanmcc/prepaid-card/models"
)

//Command - list of required command for each type transaction
type CommandInterface interface {
	AuthorisationCommand(data []string, env Env) ([]string, error)
	CaptureFundCommand(data []string, env Env) ([]string, error)
}

type Command struct{}

//AuthorisationCommand - authorise the payment
func (c *Command) AuthorisationCommand(authReq authorisationRequestBody, db *gorm.DB) error {
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
		Status:        models.STATUS_BLOCKED,
	}

	blockedTransactionRepo := models.BlockedTransactionRepository{Db: db}
	if _, err = blockedTransactionRepo.CreateBlockedTransaction(blockTransaction); err != nil {
		return err
	}

	if err := accountRepo.Update(account); err != nil {
		return err
	}

	return nil
}

//ReverseCommand - reverse the capture
func (c *Command) ReverseCommand(transactionReq *TransactionRequest, db *gorm.DB) error {
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

	if err := accountRepo.Update(account); err != nil {
		return err
	}

	if err := transactionRepo.Update(transaction); err != nil {
		return err
	}

	return nil
}
