package main

import (
	"log"
	"net/http"

	"github.com/mannanmcc/prepaid-card/handlers"
	"github.com/mannanmcc/prepaid-card/models"
)

func main() {
	db, err := models.NewDB("user:password@tcp(192.168.33.10:3306)/prepaid-card")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	env := handlers.Env{Db: db}
	router := NewRouter(env)

	//log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", router))
	log.Fatal(http.ListenAndServe(":8000", router))
}
