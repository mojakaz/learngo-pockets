package newgame

import (
	"encoding/json"
	"learngo-pockets/httpgordle/internal/api"
	"learngo-pockets/httpgordle/internal/session"
	"log"
	"net/http"
)

func Handle(w http.ResponseWriter, req *http.Request) {
	game, err := createGame()
	if err != nil {
		log.Printf("unable to create a new game: %s", err)
		http.Error(w, "failed to create a new game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	apiGame := response(game)
	err = json.NewEncoder(w).Encode(apiGame)
	if err != nil {
		// The header has already been set. Nothing much we can do here.
		log.Printf("failed to write response: %s", err)
	}
}

func createGame() (session.Game, error) {
	return session.Game{}, nil
}

func response(game session.Game) api.GameResponse {
	return api.GameResponse{}
}
