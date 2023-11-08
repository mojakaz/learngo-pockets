package guess

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"learngo-pockets/httpgordle/internal/api"
	"learngo-pockets/httpgordle/internal/session"
	"log"
	"net/http"
)

type gameFinderUpdater interface {
	Find(gameID session.GameID) (session.Game, error)
	Update(game session.Game) error
}

func Handler(db gameFinderUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, api.GameID)
		if id == "" {
			http.Error(w, "missing the id of the game", http.StatusBadRequest)
			return
		}
		// Read the request, containing the guess, from the body of the input.
		gr := api.GuessRequest{}

		game, err := guess(id, gr)
		if err != nil {
			//
		}
		_ = api.ToGameResponse(game)
		err = json.NewDecoder(r.Body).Decode(&gr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("received guess: %s", gr)
	}
}

func guess(id string, gr api.GuessRequest) (session.Game, error) {
	return session.Game{
		ID: session.GameID(id),
	}, nil
}
