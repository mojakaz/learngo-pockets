package repository

import (
	"github.com/stretchr/testify/assert"
	"learngo-pockets/httpgordle/internal/session"
	"testing"
)

func TestGameRepository_Add_Success(t *testing.T) {
	tt := map[string]struct {
		gr   *GameRepository
		game session.Game
	}{
		"nominal": {
			gr: New(),
			game: session.Game{
				ID:           "a",
				AttemptsLeft: 1,
				Guesses:      []session.Guess{},
				Status:       "Playing",
			},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := tc.gr.Add(tc.game)
			assert.NoError(t, err, "assume no error but got an error %v", err)
		})
	}
}

func TestGameRepository_Add_Failure(t *testing.T) {
	tt := map[string]struct {
		gr   *GameRepository
		game session.Game
	}{
		"key exists": {
			gr: New(),
			game: session.Game{
				ID:           "b",
				AttemptsLeft: 2,
				Guesses:      make([]session.Guess, 0),
				Status:       "Playing",
			},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			tc.gr.storage[tc.game.ID] = tc.game
			err := tc.gr.Add(tc.game)
			assert.Error(t, err, "assume an error but got no error")
		})

	}
}

func TestGameRepository_Find_Success(t *testing.T) {
	tt := map[string]struct {
		gr     *GameRepository
		gameID session.GameID
	}{
		"nominal": {
			gr:     New(),
			gameID: "exist",
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			tc.gr.storage[tc.gameID] = session.Game{}
			_, err := tc.gr.Find(tc.gameID)
			assert.NoError(t, err)
		})
	}
}

func TestGameRepository_Find_Failure(t *testing.T) {
	tt := map[string]struct {
		gr     *GameRepository
		gameID session.GameID
	}{
		"no such key": {
			gr:     New(),
			gameID: "nonexist",
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := tc.gr.Find(tc.gameID)
			assert.Error(t, err)
		})
	}
}

func TestGameRepository_Update_Success(t *testing.T) {
	tt := map[string]struct {
		gr   *GameRepository
		game session.Game
	}{
		"nominal": {
			gr: New(),
			game: session.Game{
				ID:           "a",
				AttemptsLeft: 5,
				Guesses:      make([]session.Guess, 0),
				Status:       "Playing",
			},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			tc.gr.storage[tc.game.ID] = session.Game{}
			err := tc.gr.Update(tc.game)
			assert.NoError(t, err, "expected no error, got an error: %v", err)
		})
	}
}

func TestGameRepository_Update_Failure(t *testing.T) {
	tt := map[string]struct {
		gr   *GameRepository
		game session.Game
	}{
		"nominal": {
			gr: New(),
			game: session.Game{
				ID:           "a",
				AttemptsLeft: 5,
				Guesses:      make([]session.Guess, 0),
				Status:       "Playing",
			},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := tc.gr.Update(tc.game)
			assert.Error(t, err, "expected an error, got no error")
		})
	}
}
