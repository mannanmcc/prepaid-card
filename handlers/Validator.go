package handlers

import (
	"errors"
	"net/http"
	"strconv"
)

type ValidatorInterface interface {
	Validate(r *http.Request) bool
}

//AuthorisationRequestBody - request
type authorisationRequestBody struct {
	cardNumber string
	merchantId string
	amount     float64
	reason     string
}

//Validate - validate the request
func (ar *authorisationRequestBody) Validate(r *http.Request) error {
	cardNumber := r.FormValue("cardNumber")
	if cardNumber == "" {
		return errors.New("Valid `cardNumber` must be provided`")
	}
	ar.cardNumber = cardNumber

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		return errors.New("Valid `amount` must be provided`")
	}
	ar.amount = amount

	reason := r.FormValue("reason")
	if reason == "" {
		return errors.New("argument `reason` to capture money must be provided`")
	}
	ar.reason = reason

	merchantID := r.FormValue("merchantId")
	if merchantID == "" {
		return errors.New("argument `merchantId` must be provided`")
	}

	ar.merchantId = merchantID

	return nil
}

//AuthorisationRequestBody - request
type TransactionRequest struct {
	transactionId string
	merchantId    string
	amount        float64
}

//Validate - validate the request
func (ra *TransactionRequest) Validate(r *http.Request) error {
	transactionID := r.FormValue("transactionId")
	if transactionID == "" {
		return errors.New("Parameter `transactionId` must be provided`")
	}
	ra.transactionId = transactionID

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		return errors.New("parameter `amount` must be provided`")
	}
	ra.amount = amount

	merchantID := r.FormValue("merchantId")
	if merchantID == "" {
		return errors.New("Parameter `merchantId` must be provided`")
	}

	ra.merchantId = merchantID

	return nil
}

type TopupRequest struct {
	cardNumber string
	amount     float64
}

//Validate - validate Topup request
func (topup *TopupRequest) Validate(r *http.Request) error {
	cardNumber := r.FormValue("cardNumber")
	if cardNumber == "" {
		return errors.New("Parameter `cardNumber` must be provided`")
	}
	topup.cardNumber = cardNumber

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		return errors.New("Parameter `amount` must be provided`")
	}
	topup.amount = amount

	return nil
}
