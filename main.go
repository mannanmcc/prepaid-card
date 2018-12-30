package main

import (
	"log"
	"net/http"

	"github.com/mannanmcc/prepaid-card/handlers"
	"github.com/mannanmcc/prepaid-card/models"
)

const (
	dbHost     = "localhost"
	dbUser     = "root"
	dbPassword = "password"
	dbName     = "card"
)

func main() {
	db, err := models.NewDB(dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":3306)/" + dbName + "?charset=utf8&parseTime=True")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	env := handlers.Env{Db: db}
	router := NewRouter(env)

	log.Fatal(http.ListenAndServe(":8000", router))
}
