package handlers

import (
	"net/http"
	"strconv"

	"github.com/mannanmcc/prepaid-card/models"
)

// AddNewCompany - handle add new company request
func (env Env) TopupCard(w http.ResponseWriter, r *http.Request) {
	var err error

	accountRepo := models.AccountRepository{Db: env.Db}
	accountNumber, _ := strconv.Atoi(r.FormValue("accountNumber"))
	topupAmount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

	account, err := accountRepo.FindByAccountNumber(accountNumber)
	if err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	account.Topup(topupAmount)
	err = accountRepo.UpdateAccount(account)
	if err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	JSONResponse("SUCCESS", "Your card has been top-up successfully", w)
}

// AddNewCompany - handle add new company request
func (env Env) CheckBalanceAndLoadedAmount(w http.ResponseWriter, r *http.Request) {
	JSONResponse("SUCCESS", "Your current balance is", w)
}

// AddNewCompany - handle add new company request
func (env Env) ReservedAmount(w http.ResponseWriter, r *http.Request) {
	JSONResponse("SUCCESS", "The money is reserved for the following transaction", w)
}
