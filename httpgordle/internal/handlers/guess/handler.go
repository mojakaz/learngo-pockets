package guess

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"learngo-pockets/httpgordle/internal/api"
	"learngo-pockets/httpgordle/internal/gordle"
	"learngo-pockets/httpgordle/internal/repository"
	"learngo-pockets/httpgordle/internal/session"
	"log"
	"net/http"
	"regexp"
)

type gameGuesser interface {
	Find(gameID session.GameID) (session.Game, error)
	Update(game session.Game) error
}

func Handler(db gameGuesser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, api.GameID)
		if id == "" {
			http.Error(w, "missing the id of the game", http.StatusBadRequest)
			return
		}
		// Read the request, containing the guess, from the body of the input.
		body, err := io.ReadAll(r.Body)
		guessFinderRegexp := regexp.MustCompile(`.*"guess":"(\w+)".+`)
		str := guessFinderRegexp.FindStringSubmatch(string(body))
		if len(str) != 2 {
			http.Error(w, "invalid guess", http.StatusBadRequest)
		}
		game, err := guess(id, str[1], db)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrGameNotFound):
				http.Error(w, err.Error(), http.StatusNotFound)
			case errors.Is(err, gordle.ErrInvalidGuess):
				http.Error(w, err.Error(), http.StatusBadRequest)
			case errors.Is(err, session.ErrGameOver):
				http.Error(w, err.Error(), http.StatusForbidden)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		var gr api.GameResponse
		gr = api.ToGameResponse(game)
		err = json.NewEncoder(w).Encode(gr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("received guess: %v", gr.Guesses[len(gr.Guesses)-1])
	}
}

func guess(id, guess string, db gameGuesser) (session.Game, error) {
	game, err := db.Find(session.GameID(id))
	if err != nil {
		return session.Game{}, fmt.Errorf("unable to find game: %w", err)
	}

	if game.AttemptsLeft == 0 || game.Status == session.StatusWon {
		return session.Game{}, session.ErrGameOver
	}

	feedback, err := game.Gordle.Play(guess)
	if err != nil {
		return session.Game{}, fmt.Errorf("unable to play move: %w", err)
	}

	game.Guesses = append(game.Guesses, session.Guess{
		Word:     guess,
		Feedback: feedback.String(),
	})

	game.AttemptsLeft -= 1

	switch {
	case feedback.GameWon():
		game.Status = session.StatusWon
	case game.AttemptsLeft == 0:
		fmt.Printf("😒You've lost! Game over!\n")
		game.Status = session.StatusLost
	default:
		game.Status = session.StatusPlaying
	}

	err = db.Update(game)
	if err != nil {
		return session.Game{}, fmt.Errorf("unable to save game: %w", err)
	}

	return game, nil
}
