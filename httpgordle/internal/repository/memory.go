package repository

import (
	"fmt"
	"learngo-pockets/httpgordle/internal/session"
	"log"
	"sync"
)

const (
	ErrConflictingID = apiError("gameID already exists")
	ErrGameNotFound  = apiError("cannot find game for given gameID")
)

// GameRepository holds all the current games.
type GameRepository struct {
	mutex   sync.RWMutex
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
	log.Println("Adding a game...")

	// Lock the reading and the writing of the game.
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	_, ok := gr.storage[game.ID]
	if ok {
		return fmt.Errorf("%w (%s)", ErrConflictingID, game.ID)
	}

	gr.storage[game.ID] = game

	return nil
}

// Find returns Game if the key exists in storage
func (gr *GameRepository) Find(gameID session.GameID) (session.Game, error) {
	// Lock the reading of the game.
	gr.mutex.RLock()
	defer gr.mutex.RUnlock()

	game, ok := gr.storage[gameID]
	if !ok {
		return session.Game{}, fmt.Errorf("%w (%s)", ErrGameNotFound, gameID)
	}
	return game, nil
}

// Update updates existing key value pair
func (gr *GameRepository) Update(game session.Game) error {
	// Lock the reading and the writing of the game.
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	_, ok := gr.storage[game.ID]
	if !ok {
		return fmt.Errorf("%w (%s)", ErrGameNotFound, game.ID)
	}
	gr.storage[game.ID] = game
	return nil
}
