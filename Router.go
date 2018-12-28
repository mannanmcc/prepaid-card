package main

import (
	"github.com/gorilla/mux"
	"github.com/mannanmcc/prepaid-card/handlers"
)

//NewRouter pass a request to a handler function
func NewRouter(env handlers.Env) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/card-holder/create-card", env.CreateNewCard).Methods("POST")
	r.HandleFunc("/card-holder/topup", env.TopupCard).Methods("POST")
	r.HandleFunc("/card-holder/loaded-amount", env.CheckBalanceAndLoadedAmount).Methods("POST")

	r.HandleFunc("/merchant/charge-authorise", env.AuthoriseToCharge).Methods("POST")
	r.HandleFunc("/merchant/reverse-capture", env.ReverseCapture).Methods("POST")
	r.HandleFunc("/merchant/capture-money", env.CaptureMoney).Methods("POST")
	r.HandleFunc("/merchant/refund", env.Refund).Methods("POST")
	return r
}
