package getstatus

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
	log.Printf("retrieve status of game with id: %v", id)

	apiGame := api.GameResponse{ID: id}
	err := json.NewEncoder(w).Encode(apiGame)
	if err != nil {
		// The header has already been set. Nothing much we can do here.
		log.Printf("failed to write response: %s", err)
	}
}
