package newgame

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/oklog/ulid"
	"learngo-pockets/httpgordle/internal/api"
	"learngo-pockets/httpgordle/internal/gordle"
	"learngo-pockets/httpgordle/internal/session"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//go:embed corpus/english.txt
var englishCorpus string

type gameAdder interface {
	Add(game session.Game) error
}

// Handler returns the handler for the game creation endpoint.
func Handler(db gameAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.URL.Query().Get(api.Lang)
		if len(lang) > 0 {
			// TODO create a game in the chosen language
			fmt.Println(lang)
		}
		game, err := createGame(db)

		if err != nil {
			log.Printf("unable to create a new game: %s", err)
			http.Error(w, "failed to create a new game", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		apiGame := api.ToGameResponse(game)
		err = json.NewEncoder(w).Encode(apiGame)
		if err != nil {
			// The header has already been set. Nothing much we can do here.
			log.Printf("failed to write response: %s", err)
		}
	}
}

func createGame(db gameAdder) (session.Game, error) {
	corpus, err := gordle.ReadCorpus(englishCorpus)
	if err != nil {
		return session.Game{}, fmt.Errorf("unable to read corpus: %w", err)
	}

	game, err := gordle.New(corpus)
	if err != nil {
		return session.Game{}, fmt.Errorf("failed to create a new gordle game")
	}
	t := time.Unix(1000000, 0)
	u, err := ulid.New(ulid.Timestamp(t), ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0))
	if err != nil {
		return session.Game{}, fmt.Errorf("failed to create a new ID")
	}
	g := session.Game{
		ID:           session.GameID(u.String()),
		Gordle:       *game,
		AttemptsLeft: byte(game.ShowLength()),
		Guesses:      []session.Guess{},
		Status:       session.StatusPlaying,
		Solution:     game.ShowAnswer(),
	}

	err = db.Add(g)
	if err != nil {
		return session.Game{}, fmt.Errorf("failed to save the new game: %w", err)
	}

	return g, nil
}
