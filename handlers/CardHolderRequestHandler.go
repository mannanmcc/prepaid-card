package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/mannanmcc/prepaid-card/models"
)

//AccountDetails - return structure type
type AccountDetails struct {
	CurrentBalance float64
	BlockedFunds   []BlockedFund
	RecentToptops  []RecentTopup
}

type BlockedFund struct {
	Amount         float64
	TransactionRef string
	Reason         string
	BlockedAt      time.Time
}

type RecentTopup struct {
	Amount  float64
	TopupAt time.Time
}

// TopupCard - handle add new company request
func (env Env) TopupCard(w http.ResponseWriter, r *http.Request) {
	cardNumber := r.FormValue("cardNumber")
	topupAmount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	command := CardHolderCommand{}
	if err := command.Toptup(cardNumber, topupAmount, env.Db); err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	topup := &models.Topup{
		Amount:     topupAmount,
		CardNumber: cardNumber,
		TopupAt:    time.Now(),
	}

	topupRepo := models.TopupRepository{Db: env.Db}
	if err := topupRepo.Store(topup); err != nil {
		JSONResponse("FAILED", err.Error(), w)
		return
	}

	JSONResponse("SUCCESS", "Your card has been top-up successfully", w)
}

// CheckBalanceAndLoadedAmount - handle add new company request
func (env Env) CheckBalanceAndLoadedAmount(w http.ResponseWriter, r *http.Request) {
	cardNumber := r.FormValue("cardNumber")
	accountRepo := models.AccountRepository{Db: env.Db}
	account, err := accountRepo.FindByCardNumber(cardNumber)
	if err != nil {
		HandleFAILEDResponse(err.Error(), w)
	}

	blockedTransRepo := models.BlockedTransactionRepository{Db: env.Db}
	transactions := blockedTransRepo.FindAllBlockedTransactions(cardNumber)

	accountDetails := &AccountDetails{
		CurrentBalance: account.Balance,
	}

	accountDetails.BlockedFunds = make([]BlockedFund, 0)
	for _, transaction := range transactions {
		blockedFund := BlockedFund{
			Amount:         transaction.Balance,
			TransactionRef: transaction.TransactionID,
			Reason:         transaction.Reason,
			BlockedAt:      transaction.BlockedAt,
		}

		accountDetails.BlockedFunds = append(accountDetails.BlockedFunds, blockedFund)
	}

	topupRepo := models.TopupRepository{Db: env.Db}
	topups := topupRepo.FindAllTopups(cardNumber)

	accountDetails.RecentToptops = make([]RecentTopup, 0)
	for _, topup := range topups {
		recentTopup := RecentTopup{
			Amount:  topup.Amount,
			TopupAt: topup.TopupAt,
		}

		accountDetails.RecentToptops = append(accountDetails.RecentToptops, recentTopup)
	}

	json.NewEncoder(w).Encode(accountDetails)
	//	JSONResponse("SUCCESS", "Your current balance is", w)
}

// ReservedAmount - handle add new company request
func (env Env) ReservedAmount(w http.ResponseWriter, r *http.Request) {
	JSONResponse("SUCCESS", "The money is reserved for the following transaction", w)
}
