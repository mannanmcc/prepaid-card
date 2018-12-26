package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/mannanmcc/prepaid-card/models"
)

const RANDOM_KEY_LENGTH = 16

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

// AuthoriseToCharge - authorise to capture a fund
func (env Env) AuthoriseToCharge(w http.ResponseWriter, r *http.Request) {
	cardNumber := r.FormValue("cardNumber")
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	reason := r.FormValue("reason")
	merchantID := r.FormValue("merchantId")

	if err := handleMerchantVerification(merchantID, env.Db); err != nil {
		HandleFAILEDResponse(err.Error(), w)
		return
	}

	command := Command{}
	if err := command.AuthorisationCommand(merchantID, cardNumber, amount, reason, env.Db); err != nil {
		HandleFAILEDResponse(err.Error(), w)
	}

	JSONResponse("SUCCESS", "This authorisation request is approved.", w)
}

// ReverseCapture - reverse the transaction and the amount can not be charge to the again
func (env Env) ReverseCapture(w http.ResponseWriter, r *http.Request) {
	transactionID := r.FormValue("transactionId")
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

	command := Command{}
	command.ReverseCommand(transactionID, amount, env.Db)
	JSONResponse("SUCCESS", "The transaction with ref XXXXX is reversed and can not be charge back to the card", w)
}

// CaptureMoney - captureMoney captures amount for merchant
func (env Env) CaptureMoney(w http.ResponseWriter, r *http.Request) {
	transactionID := r.FormValue("transactionId")
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

	command := CaptureCommand{}
	if err := command.CaptureFund(transactionID, amount, env.Db); err != nil {
		HandleFAILEDResponse(err.Error(), w)
		return
	}

	JSONResponse("SUCCESS", "The amount is capture and sent to merchant", w)
}

// Refund - handle to refund by merchant
func (env Env) Refund(w http.ResponseWriter, r *http.Request) {
	transactionID := r.FormValue("transactionId")
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	command := CaptureCommand{}

	if err := command.Refund(transactionID, amount, env.Db); err != nil {
		HandleFAILEDResponse(err.Error(), w)
		return
	}

	JSONResponse("SUCCESS", "The transaction with ref XXXXX has been refund", w)
}
