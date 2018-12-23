package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

//TransactionRepositoryInterface - interface for model TransactionRepository
type TransactionRepositoryInterface interface {
	createTransaction(transaction Transaction) (Transaction, error)
}

//TransactionRepository - type TransactionRepository
type TransactionRepository struct {
	Db *gorm.DB
}

//CreateTransaction - record a block transaction
func (repo *TransactionRepository) Create(transaction Transaction) (Transaction, error) {
	id := repo.Db.Save(&transaction)

	if id == nil {
		return transaction, errors.New("account saving failed")
	}

	return transaction, nil
}

//FindByTransactionId - find the block transaction by transaction id
func (repo *TransactionRepository) FindByTransactionID(transationID string) (*Transaction, error) {
	var transaction Transaction
	res := repo.Db.Find(&transaction, &Transaction{TransactionID: transationID, Status: STATUS_CAPTURED})

	if res.RecordNotFound() {
		return nil, fmt.Errorf("No transaction found to refund with id: %s", transationID)
	}

	return &transaction, nil
}

//Update - update the blocked transaction
func (repo *TransactionRepository) Update(transaction *Transaction) error {
	id := repo.Db.Save(&transaction)
	if id != nil {
		return errors.New("block transaction saving failed")
	}

	return nil
}
