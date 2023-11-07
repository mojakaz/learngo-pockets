package api

import "net/http"

const (
	NewGameRoute   = "/games"
	NewGameMethod  = http.MethodPost
	GameID         = "id"
	GetStatusRoute = "/games/{" + GameID + "}"
)

type GameResponse struct {
	ID           string  `json:"id"`
	AttemptsLeft byte    `json:"attempts_left"`
	Guesses      []Guess `json:"guesses"`
	WordLength   byte    `json:"word_length"`
	Solution     string  `json:"solution,omitempty"`
	Status       string  `json:"status"`
}

type Guess struct {
	Word     string `json:"word"`
	Feedback string `json:"feedback"`
}
