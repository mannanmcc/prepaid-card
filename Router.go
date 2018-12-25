package main

import (
	"github.com/gorilla/mux"
	"github.com/mannanmcc/prepaid-card/handlers"
)

//NewRouter points a request to a handler function
func NewRouter(env handlers.Env) *mux.Router {
	r := mux.NewRouter()

	/*todo - all method name should be start with uppercase or lowercase*/
	r.HandleFunc("/merchant/charge-authorise", env.AuthoriseToCharge).Methods("POST")
	r.HandleFunc("/merchant/capture-money", env.CaptureMoney).Methods("POST")
	r.HandleFunc("/merchant/reverse-capture", env.ReverseCapture).Methods("POST")
	r.HandleFunc("/merchant/refund", env.Refund).Methods("POST")

	r.HandleFunc("/card-holder/topup", env.TopupCard).Methods("POST")
	r.HandleFunc("/card-holder/loaded-amount", env.CheckBalanceAndLoadedAmount).Methods("GET")
	r.HandleFunc("/card-holder/amount-reserved", env.ReservedAmount).Methods("GET")

	return r
}
