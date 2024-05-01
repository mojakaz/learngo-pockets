package getstatus

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"learngo-pockets/httpgordle/internal/api"
	"learngo-pockets/httpgordle/internal/repository"
	"learngo-pockets/httpgordle/internal/session"
	"log"
	"net/http"
)

type gameFinder interface {
	Find(gameID session.GameID) (session.Game, error)
}

func Handler(db gameFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, api.GameID)
		if id == "" {
			http.Error(w, "missing the id of the game", http.StatusBadRequest)
			return
		}
		log.Printf("retrieve status of game with id: %v", id)
		game, err := db.Find(session.GameID(id))
		if err != nil {
			if errors.Is(err, repository.ErrGameNotFound) {
				http.Error(w, "this game does not exist", http.StatusNotFound)
				return
			}

			log.Printf("cannot fetch game %s: %s", id, err)
			http.Error(w, "failed to fetch game", http.StatusInternalServerError)
			return
		}
		apiGame := api.ToGameResponse(game)
		err = json.NewEncoder(w).Encode(apiGame)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
