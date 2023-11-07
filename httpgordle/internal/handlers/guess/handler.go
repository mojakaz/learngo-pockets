package guess

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"learngo-pockets/httpgordle/internal/api"
	"log"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, api.GameID)
	if id == "" {
		http.Error(w, "missing the id of the game", http.StatusBadRequest)
		return
	}
	// Read the request, containing the guess, from the body of the input.
	guess := api.GuessRequest{}
	err := json.NewDecoder(r.Body).Decode(&guess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("received guess: %s", guess)
}
