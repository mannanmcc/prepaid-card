package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/mannanmcc/prepaid-card/models"
)

const RANDOM_KEY_LENGTH = 16

// AuthoriseToCharge - authorise to capture a fund
func (env Env) AuthoriseToCharge(w http.ResponseWriter, r *http.Request) {
	var err error

	accountRepo := models.AccountRepository{Db: env.Db}
	//todo- should we find by account number or card number???
	accountNumber, _ := strconv.Atoi(r.FormValue("accountNumber"))
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	merchantID := r.FormValue("merchantId")
	reason := r.FormValue("reason")

	account, err := accountRepo.FindByAccountNumber(accountNumber)
	if err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	//check if the account has enough money
	_, err = account.AuthoriseAmount(amount)
	if err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	transactionID := models.GenerateTransactionId(RANDOM_KEY_LENGTH)
	currentTime := time.Now()
	blockTransaction := models.BlockedTransaction{
		TransactionID: transactionID,
		AccountNumber: account.AccountNumber,
		Amount:        amount,
		MerchantID:    merchantID,
		Reason:        reason,
		BlockedAt:     currentTime,
		Status:        models.STATUS_BLOCKED,
	}

	blockedTransactionRepo := models.BlockedTransactionRepository{Db: env.Db}
	if _, err = blockedTransactionRepo.CreateBlockedTransaction(blockTransaction); err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	//update the balance
	accountRepo.UpdateAccount(account)
	JSONResponseWithTransaction(transactionID, "SUCCESS", "This authorisation request is approved.", w)
}

// CaptureMoney - captureMoney captures amount for merchant
func (env Env) CaptureMoney(w http.ResponseWriter, r *http.Request) {
	transactionID := r.FormValue("transactionId")
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	blockTransactionRepo := models.BlockedTransactionRepository{Db: env.Db}
	blockedTransaction, err := blockTransactionRepo.FindByTransactionID(transactionID)
	if err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	if err := blockedTransaction.CaptureFund(amount); err != nil {
		JSONResponseWithTransaction(transactionID, "FAILED", err.Error(), w)
		return
	}

	transactionID = models.GenerateTransactionId(RANDOM_KEY_LENGTH)
	transactionRepo := models.TransactionRepository{Db: env.Db}
	currentTime := time.Now()

	transaction := models.Transaction{
		TransactionID:        transactionID,
		BlockedTransactionID: blockedTransaction.TransactionID,
		MerchantID:           blockedTransaction.MerchantID,
		Amount:               amount,
		Status:               models.STATUS_CAPTURED,
		AccountNumber:        blockedTransaction.AccountNumber,
		CapturedAt:           currentTime,
	}

	if _, err := transactionRepo.Create(transaction); err != nil {
		JSONResponse("FAILED", err.Error(), w)
	}

	blockTransactionRepo.Update(blockedTransaction)
	JSONResponse("SUCCESS", "The amount is capture and sent to merchant", w)
}

// ReverseCapture - reverse the transaction and the amount can not be charge to the again
func (env Env) ReverseCapture(w http.ResponseWriter, r *http.Request) {
	JSONResponse("SUCCESS", "The transaction with ref XXXXX is reversed and can not be charge back to the card", w)
}

// Refund - handle to refund by merchant
func (env Env) Refund(w http.ResponseWriter, r *http.Request) {
	transactionID := r.FormValue("transactionId")
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	transactionRepo := models.TransactionRepository{Db: env.Db}
	transaction, err := transactionRepo.FindByTransactionID(transactionID)
	if err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	accountRepo := models.AccountRepository{Db: env.Db}
	account, err := accountRepo.FindByAccountNumber(transaction.AccountNumber)
	if err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	err = account.Refund(transaction, amount)
	if err := accountRepo.UpdateAccount(account); err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	err = transactionRepo.Update(transaction)

	JSONResponse("SUCCESS", "The transaction with ref XXXXX has been refund", w)
}
