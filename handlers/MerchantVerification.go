package handlers

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/mannanmcc/prepaid-card/models"
)

func handleMerchantVerification(merchantID string, db *gorm.DB) error {
	var merchant models.Merchant
	merchantRepo := models.MerchantRepository{Db: db}
	merchant, err := merchantRepo.FindByMerchantID(merchantID)
	if err != nil {
		return fmt.Errorf("no merchant is found with merchant id: %s", merchantID)
	}

	if merchant.Status != models.MERCHANT_STATUS_ACTIVE {
		return fmt.Errorf("Merchant %s is not active", merchantID)
	}

	return nil
}

func isBlockedTransactionBelongsToMerchant(transactionID string, merchantID string, db *gorm.DB) error {
	transactionRepo := models.BlockedTransactionRepository{Db: db}
	transaction, err := transactionRepo.FindByTransactionID(transactionID)
	if err != nil {
		return fmt.Errorf("no transaction is found with transaction id: %s", transactionID)
	}

	if transaction.MerchantID != merchantID {
		return fmt.Errorf("This transaction is not belongs to this merchant")
	}

	return nil
}

func isTransactionBelongsToMerchant(transactionID string, merchantID string, db *gorm.DB) error {
	transactionRepo := models.TransactionRepository{Db: db}
	transaction, err := transactionRepo.FindByTransactionID(transactionID)
	if err != nil {
		return fmt.Errorf("no transaction is found with transaction id: %s", transactionID)
	}

	if transaction.MerchantID != merchantID {
		return fmt.Errorf("This transaction is not belongs to this merchant")
	}

	return nil
}
