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

//JSONResponse builds up the response object and encode
func JSONResponse(status string, msg string, w http.ResponseWriter) {
	response := Response{
		Status:  status,
		Message: msg,
	}

	json.NewEncoder(w).Encode(response)
}
