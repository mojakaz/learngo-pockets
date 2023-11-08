package handlers

import (
	"github.com/go-chi/chi/v5"
	"learngo-pockets/httpgordle/internal/api"
	"learngo-pockets/httpgordle/internal/handlers/getstatus"
	"learngo-pockets/httpgordle/internal/handlers/guess"
	"learngo-pockets/httpgordle/internal/handlers/newgame"
	"learngo-pockets/httpgordle/internal/repository"
)

func NewRouter(db *repository.GameRepository) chi.Router {
	r := chi.NewRouter()

	r.Post(api.NewGameRoute, newgame.Handler(db))
	r.Get(api.GetStatusRoute, getstatus.Handler(db))
	r.Put(api.GuessRoute, guess.Handler(db))

	return r
}
