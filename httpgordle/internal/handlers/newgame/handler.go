package newgame

import (
	"encoding/json"
	"learngo-pockets/httpgordle/internal/api"
	"learngo-pockets/httpgordle/internal/session"
	"log"
	"net/http"
)

type gameAdder interface {
	Add(game session.Game) error
}

// Handler returns the handler for the game creation endpoint.
func Handler(db gameAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		game, err := createGame(db)

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
}

func createGame(db gameAdder) (session.Game, error) {
	return session.Game{}, nil
}

func response(game session.Game) api.GameResponse {
	return api.GameResponse{}
}
