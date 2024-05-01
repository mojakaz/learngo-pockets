package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"learngo-pockets/httpgordle/internal/api"
	"learngo-pockets/httpgordle/internal/handlers/deletegame"
	"learngo-pockets/httpgordle/internal/handlers/getstatus"
	"learngo-pockets/httpgordle/internal/handlers/guess"
	"learngo-pockets/httpgordle/internal/handlers/newgame"
	"learngo-pockets/httpgordle/internal/repository"
)

func NewRouter(db *repository.GameRepository) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Throttle(10))
	r.Post(api.NewGameRoute, newgame.Handler(db))
	r.Get(api.GetStatusRoute, getstatus.Handler(db))
	r.Put(api.GuessRoute, guess.Handler(db))
	r.Delete(api.DeleteGameRoute, deletegame.Handler(db))

	return r
}
