package handlers

import (
	"net/http"
	"strconv"

	"github.com/mannanmcc/prepaid-card/models"
)

// AuthoriseToCharge - authorise to capture a fund
func (env Env) AuthoriseToCharge(w http.ResponseWriter, r *http.Request) {
	var err error

	accountRepo := models.AccountRepository{Db: env.Db}
	//todo- should we find by account number or card number???
	accountNumber, _ := strconv.Atoi(r.FormValue("accountNumber"))
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

	account, err := accountRepo.FindByAccountNumber(accountNumber)
	if err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	//check if the account has amount money
	_, err = account.AuthoriseAmount(amount)
	if err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	JSONResponse("SUCCESS", "This authorisation request is approved", w)
}

// CaptureMoney - captureMoney captures amount for merchant
func (env Env) CaptureMoney(w http.ResponseWriter, r *http.Request) {
	JSONResponse("SUCCESS", "The amount is capture and will be send to merchant very soon", w)
}

// ReverseCapture - reverse the transaction and the amount can not be charge to the again
func (env Env) ReverseCapture(w http.ResponseWriter, r *http.Request) {
	JSONResponse("SUCCESS", "The transaction with ref XXXXX is reversed and can not be charge back to the card", w)
}

// Refund - handle to refund by merchant
func (env Env) Refund(w http.ResponseWriter, r *http.Request) {
	JSONResponse("SUCCESS", "The transaction with ref XXXXX has been refund", w)
}
