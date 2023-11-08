package repository

import (
	"fmt"
	"learngo-pockets/httpgordle/internal/session"
)

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

// Add inserts for the first time a game in memory.
func (gr *GameRepository) Add(game session.Game) error {
	_, ok := gr.storage[game.ID]
	if ok {
		return fmt.Errorf("gameID %s already exists", game.ID)
	}

	gr.storage[game.ID] = game

	return nil
}

// Find returns Game if the key exists in storage
func (gr *GameRepository) Find(gameID session.GameID) (session.Game, error) {
	game, ok := gr.storage[gameID]
	if !ok {
		return session.Game{}, fmt.Errorf("cannot find game for given gameID %s", gameID)
	}
	return game, nil
}

// Update updates existing key value pair
func (gr *GameRepository) Update(game session.Game) error {
	_, ok := gr.storage[game.ID]
	if !ok {
		return fmt.Errorf("cannnot find game for given gameID %s", game.ID)
	}
	gr.storage[game.ID] = game
	return nil
}
