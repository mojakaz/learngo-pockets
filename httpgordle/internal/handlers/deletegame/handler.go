package deletegame

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"learngo-pockets/httpgordle/internal/api"
	"learngo-pockets/httpgordle/internal/session"
	"log"
	"net/http"
)

type gameDeleter interface {
	Delete(gameID session.GameID) error
}

func Handler(db gameDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, api.GameID)
		if id == "" {
			http.Error(w, "missing the id of the game", http.StatusBadRequest)
			return
		}
		log.Printf("delete game with id: %v", id)
		err := deleteGame(id, db)

		if err != nil {
			log.Println(err.Error())
			http.Error(w, "failed to delete game", http.StatusInternalServerError)
			return
		}
	}
}

func deleteGame(id string, db gameDeleter) error {
	err := db.Delete(session.GameID(id))
	if err != nil {
		return fmt.Errorf("unable to delete game: %w", err)
	}
	return err
}
