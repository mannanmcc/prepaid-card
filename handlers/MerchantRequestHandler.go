package handlers

import (
	"fmt"
	"net/http"
)

const RANDOM_KEY_LENGTH = 16

// AuthoriseToCharge - verify the card balance and get approval to charge
func (env Env) AuthoriseToCharge(w http.ResponseWriter, r *http.Request) {
	authorisationRequest := authorisationRequestBody{}
	if err := authorisationRequest.Validate(r); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	if err := handleMerchantVerification(authorisationRequest.merchantId, env.Db); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	command := Command{}
	if err := command.AuthorisationCommand(authorisationRequest, env.Db); err != nil {
		HandleFailedResponse(err.Error(), w)
	}

	JSONResponse("SUCCESS", "This authorisation request is approved.", w)
}

// ReverseCapture - reverse the transaction and the amount can not be charge to the again
func (env Env) ReverseCapture(w http.ResponseWriter, r *http.Request) {
	transReq := &TransactionRequest{}
	if err := transReq.Validate(r); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	if err := handleMerchantVerification(transReq.merchantId, env.Db); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	if err := isBlockedTransactionBelongsToMerchant(transReq.transactionId, transReq.merchantId, env.Db); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	command := Command{}
	if err := command.ReverseCommand(transReq, env.Db); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	JSONResponse("SUCCESS", fmt.Sprintf("The transaction with ref %s is reversed and can not be charge back to the card", transReq.transactionId), w)
}

// CaptureMoney - captureMoney captures amount for merchant
func (env Env) CaptureMoney(w http.ResponseWriter, r *http.Request) {
	captureRequest := &TransactionRequest{}
	if err := captureRequest.Validate(r); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	if err := handleMerchantVerification(captureRequest.merchantId, env.Db); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	if err := isBlockedTransactionBelongsToMerchant(captureRequest.transactionId, captureRequest.merchantId, env.Db); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	command := CaptureCommand{}
	if err := command.CaptureFund(captureRequest, env.Db); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	JSONResponse("SUCCESS", "The amount is capture and sent to merchant", w)
}

// Refund - handle to refund by merchant
func (env Env) Refund(w http.ResponseWriter, r *http.Request) {
	refundReq := &TransactionRequest{}

	if err := refundReq.Validate(r); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	if err := handleMerchantVerification(refundReq.merchantId, env.Db); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	if err := isTransactionBelongsToMerchant(refundReq.transactionId, refundReq.merchantId, env.Db); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	command := CaptureCommand{}
	if err := command.Refund(refundReq, env.Db); err != nil {
		HandleFailedResponse(err.Error(), w)
		return
	}

	JSONResponse("SUCCESS", fmt.Sprintf("The transaction with ref %s has been refund", refundReq.transactionId), w)
}
