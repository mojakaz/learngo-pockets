package main

import (
	"learngo-pockets/httpgordle/internal/handlers"
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8080", handlers.NewRouter())
	if err != nil {
		log.Fatal(err)
	}
}
