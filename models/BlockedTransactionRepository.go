package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

//BlockedTransactionRepositoryInterface - interface for model BlockedTransactionRepository
type BlockedTransactionRepositoryInterface interface {
	TopupAccount(account Account) (Account, error)
	FindByAccountNumber(accountNumber int) (*Account, error)
}

//BlockedTransactionRepository - type BlockedTransactionRepository
type BlockedTransactionRepository struct {
	Db *gorm.DB
}

//CreateBlockedTransaction - record a block transaction
func (repo *BlockedTransactionRepository) CreateBlockedTransaction(transaction BlockedTransaction) (BlockedTransaction, error) {
	id := repo.Db.Save(&transaction)

	if id == nil {
		return transaction, errors.New("account saving failed")
	}

	return transaction, nil
}

//FindByTransactionId - find the block transaction by transaction id
func (repo *BlockedTransactionRepository) FindByTransactionID(transationID string) (*BlockedTransaction, error) {
	var blockedTransaction BlockedTransaction
	res := repo.Db.Find(&blockedTransaction, &BlockedTransaction{TransactionID: transationID, Status: STATUS_BLOCKED})

	if res.RecordNotFound() {
		return nil, fmt.Errorf("No blocked transaction found to capture amount from with id: %s", transationID)
	}

	return &blockedTransaction, nil
}

//Update - update the blocked transaction
func (repo *BlockedTransactionRepository) Update(transaction *BlockedTransaction) error {
	id := repo.Db.Save(&transaction)
	if id != nil {
		return errors.New("block transaction saving failed")
	}

	return nil
}
