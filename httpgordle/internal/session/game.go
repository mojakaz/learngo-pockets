package session

import (
	"errors"
	"learngo-pockets/httpgordle/internal/gordle"
)

const (
	StatusPlaying = "Playing"
	StatusWon     = "Won"
	StatusLost    = "Lost"
)

// ErrGameOver is returned when a play is made but the game is over.
var ErrGameOver = errors.New("game over")

// Game contains the information about a game.
type Game struct {
	ID           GameID
	Gordle       gordle.Game
	AttemptsLeft byte
	Guesses      []Guess
	Status       Status
	Solution     string
}

// A GameID represents the ID of a game.
type GameID string

func (g GameID) String() string {
	return string(g)
}

// A Guess is a pair of a word (submitted by the player) and its feedback (provided by Gordle).
type Guess struct {
	Word     string
	Feedback string
}

// Status is the current status of the game and tells what operations can be made on it.
type Status string
