package repository

import "learngo-pockets/httpgordle/internal/session"

// GameRepository holds all the current games.
type GameRepository struct {
	storage map[session.GameID]session.Game
}

// New creates an empty game repository.
func New() *GameRepository {
	return &GameRepository{
		storage: make(map[session.GameID]session.Game),
	}
}
