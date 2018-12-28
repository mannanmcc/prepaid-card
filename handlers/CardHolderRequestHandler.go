package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mannanmcc/prepaid-card/models"
)

//CreateCardSuccessResponse - response structure for create new card
type CreateCardSuccessResponse struct {
	CardNumber string
	Message    string
}

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
	topupRequest := &TopupRequest{}
	if err := topupRequest.Validate(r); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}
	command := CardHolderCommand{}
	if err := command.Toptup(topupRequest, env.Db); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	topup := &models.Topup{
		Amount:     topupRequest.amount,
		CardNumber: topupRequest.cardNumber,
		TopupAt:    time.Now(),
	}

	topupRepo := models.TopupRepository{Db: env.Db}
	if err := topupRepo.Store(topup); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	JSONResponse("SUCCESS", "Your card has been top-up successfully", w)
}

// CreateNewCard - create a new card
func (env Env) CreateNewCard(w http.ResponseWriter, r *http.Request) {
	createCardRequest := &CreateCardRequest{}
	if err := createCardRequest.Validate(r); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	command := CardHolderCommand{}
	account, err := command.CreateCard(createCardRequest, env.Db)
	if err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	response := &CreateCardSuccessResponse{
		CardNumber: account.CardNumber,
		Message:    "Welcome!, new card has been created",
	}

	json.NewEncoder(w).Encode(response)
}

// CheckBalanceAndLoadedAmount - handle add new company request
func (env Env) CheckBalanceAndLoadedAmount(w http.ResponseWriter, r *http.Request) {
	cardNumber := r.FormValue("cardNumber")
	accountRepo := models.AccountRepository{Db: env.Db}
	account, err := accountRepo.FindByCardNumber(cardNumber)
	if err != nil {
		HandleFailedResponse(err.Error(), w)
	}

	blockedTransRepo := models.BlockedTransactionRepository{Db: env.Db}
	transactions := blockedTransRepo.FindAllBlockedTransactions(cardNumber)

	accountDetails := &AccountDetails{
		CurrentBalance: account.Balance,
	}

	accountDetails.BlockedFunds = make([]BlockedFund, 0)
	for _, transaction := range transactions {
		if transaction.Status == models.StatusBlocked && transaction.Balance > 0 {
			blockedFund := BlockedFund{
				Amount:         transaction.Balance,
				TransactionRef: transaction.TransactionID,
				Reason:         transaction.Reason,
				BlockedAt:      transaction.BlockedAt,
			}
			accountDetails.BlockedFunds = append(accountDetails.BlockedFunds, blockedFund)
		}
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
}
