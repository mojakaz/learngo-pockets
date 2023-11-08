package main

import (
	"learngo-pockets/httpgordle/internal/handlers"
	"learngo-pockets/httpgordle/internal/repository"
	"log"
	"net/http"
)

func main() {
	db := repository.New()
	err := http.ListenAndServe(":8080", handlers.NewRouter(db))
	if err != nil {
		log.Fatal(err)
	}
}
