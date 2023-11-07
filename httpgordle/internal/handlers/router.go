package handlers

import (
	"github.com/go-chi/chi/v5"
	"learngo-pockets/httpgordle/internal/api"
	"learngo-pockets/httpgordle/internal/handlers/getstatus"
	"learngo-pockets/httpgordle/internal/handlers/guess"
	"learngo-pockets/httpgordle/internal/handlers/newgame"
	"net/http"
)

func Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(api.NewGameRoute, newgame.Handle)
	return mux
}

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Post(api.NewGameRoute, newgame.Handle)
	r.Get(api.GetStatusRoute, getstatus.Handle)
	r.Put(api.GuessRoute, guess.Handle)

	return r
}
