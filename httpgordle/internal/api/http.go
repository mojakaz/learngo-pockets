package api

import "net/http"

const (
	// POST /games - creates a new game and returns its ID.
	NewGameRoute  = "/games"
	NewGameMethod = http.MethodPost
)

type GameResponse struct {
	Id           string  `json:"id"`
	AttemptsLeft int     `json:"attempts_left"`
	Guesses      []Guess `json:"guesses"`
	WordLength   int     `json:"word_length"`
	Solution     string  `json:"solution,omitempty"`
	Status       string  `json:"status"`
}

type Guess struct {
	Word     string `json:"word"`
	Feedback string `json:"feedback"`
}
