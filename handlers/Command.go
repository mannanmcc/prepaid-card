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

func (c *Command) AuthorisationCommand(merchantID string, cardNumber string, amount float64, reason string, db *gorm.DB) error {
	accountRepo := models.AccountRepository{Db: db}
	account, err := accountRepo.FindByCardNumber(cardNumber)
	if err != nil {
		return err
	}

	_, err = account.AuthoriseAmount(amount)
	if err != nil {
		return err
	}

	transactionID := models.GenerateTransactionId(RANDOM_KEY_LENGTH)
	blockTransaction := models.BlockedTransaction{
		TransactionID: transactionID,
		CardNumber:    account.CardNumber,
		Amount:        amount,
		Balance:       amount,
		MerchantID:    merchantID,
		Reason:        reason,
		BlockedAt:     time.Now(),
		Status:        models.STATUS_BLOCKED,
	}

	blockedTransactionRepo := models.BlockedTransactionRepository{Db: db}
	if _, err = blockedTransactionRepo.CreateBlockedTransaction(blockTransaction); err != nil {
		return err
	}

	if err := accountRepo.UpdateAccount(account); err != nil {
		return err
	}

	return nil
}

func (c *Command) ReverseCommand(transactionID string, amount float64, db *gorm.DB) error {
	transactionRepo := models.BlockedTransactionRepository{Db: db}
	transaction, err := transactionRepo.FindByTransactionID(transactionID)
	if err != nil {
		return err
	}

	accountRepo := models.AccountRepository{Db: db}
	account, err := accountRepo.FindByCardNumber(transaction.CardNumber)
	if err != nil {
		return err
	}

	if err = transaction.Reverse(amount); err != nil {
		return err
	}

	account.Topup(amount)

	if err := accountRepo.UpdateAccount(account); err != nil {
		return err
	}

	if err := transactionRepo.Update(transaction); err != nil {
		return err
	}

	return nil
}
