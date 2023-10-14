package newgame

import (
	"learngo-pockets/httpgordle/internal/api"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != api.NewGameMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("Creating a new game"))
}
