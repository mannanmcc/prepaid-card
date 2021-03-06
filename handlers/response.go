package handlers

import (
	"encoding/json"
	"net/http"
)

//Response - represent response type
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

//BlockedTransactionResponse - response for blocking an amount on card
type BlockedTransactionResponse struct {
	TransactionID string
	Status        string
	Message       string
}

//JSONResponse builds up the response object and encode
func JSONResponseWithTransaction(transactionId string, status string, msg string, w http.ResponseWriter) {
	response := BlockedTransactionResponse{
		TransactionID: transactionId,
		Status:        status,
		Message:       msg,
	}

	json.NewEncoder(w).Encode(response)
}

//JSONResponse builds up the success response object and send
func handleSuccessResponse(msg string, w http.ResponseWriter) {
	response := Response{
		Status:  "SUCCESS",
		Message: msg,
	}

	json.NewEncoder(w).Encode(response)
}

func handleFailedResponse(msg string, w http.ResponseWriter) {
	response := Response{
		Status:  "FAILED",
		Message: msg,
	}

	json.NewEncoder(w).Encode(response)
}
