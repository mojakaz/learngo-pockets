package api

import (
	"github.com/stretchr/testify/assert"
	"learngo-pockets/httpgordle/internal/gordle"
	"learngo-pockets/httpgordle/internal/session"
	"testing"
)

func TestToGameResponse(t *testing.T) {
	tt := map[string]struct {
		game     session.Game
		expected GameResponse
		err      error
	}{
		"new": {
			game: session.Game{
				ID:           "a",
				Gordle:       gordle.Game{},
				AttemptsLeft: 5,
				Guesses:      nil,
				Status:       "Playing",
				Solution:     "hello",
			},
			expected: GameResponse{
				ID:           "a",
				AttemptsLeft: 5,
				Guesses:      []session.Guess{},
				WordLength:   5,
				Solution:     "",
				Status:       "Playing",
			},
			err: nil,
		},
		"two guesses": {
			game: session.Game{
				ID:           "b",
				Gordle:       gordle.Game{},
				AttemptsLeft: 3,
				Guesses: []session.Guess{{
					Word:     "abc",
					Feedback: "",
				}, {
					Word:     "def",
					Feedback: "",
				}},
				Status:   "Playing",
				Solution: "xyz",
			},
			expected: GameResponse{
				ID:           "b",
				AttemptsLeft: 3,
				Guesses: []session.Guess{{
					Word:     "abc",
					Feedback: "",
				}, {
					Word:     "def",
					Feedback: "",
				}},
				WordLength: 3,
				Solution:   "",
				Status:     "Playing",
			},
			err: nil,
		},
		"zero attempts left": {
			game: session.Game{
				ID:           "c",
				Gordle:       gordle.Game{},
				AttemptsLeft: 0,
				Guesses:      nil,
				Status:       "LOST",
				Solution:     "NIGHT",
			},
			expected: GameResponse{
				ID:           "c",
				AttemptsLeft: 0,
				Guesses:      make([]session.Guess, 0),
				WordLength:   5,
				Solution:     "NIGHT",
				Status:       "LOST",
			},
			err: nil,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := ToGameResponse(tc.game)
			assert.Equalf(t, got, tc.expected, "expected %v, got %v", tc.expected, got)
		})
	}
}
