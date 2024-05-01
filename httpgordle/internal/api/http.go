package api

import (
	"learngo-pockets/httpgordle/internal/session"
	"net/http"
)

const (
	NewGameRoute    = "/games"
	NewGameMethod   = http.MethodPost
	GameID          = "id"
	GetStatusRoute  = "/games/{" + GameID + "}"
	DeleteGameRoute = "/games/{" + GameID + "}"
	Lang            = "lang"
)

type GameResponse struct {
	ID           string          `json:"id"`
	AttemptsLeft byte            `json:"attempts_left"`
	Guesses      []session.Guess `json:"guesses"`
	WordLength   byte            `json:"word_length"`
	Solution     string          `json:"solution,omitempty"`
	Status       string          `json:"status"`
}
