package newgame

import (
	"github.com/stretchr/testify/assert"
	"learngo-pockets/httpgordle/internal/session"
	"testing"
)

func TestCreateGame(t *testing.T) {
	tt := map[string]struct {
		db  gameAdder
		err error
	}{
		"nominal": {
			db:  gameAdderStub2{},
			err: nil,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			game, err := createGame(tc.db)
			if err != nil {
				t.Fatalf("expected no error, but got %v", err)
			}
			// Assert game.ID consists of digits and letters(uppercase and/or lowercase)
			assert.Regexp(t, `\w+`, game.ID)
			assert.Equal(t, uint8(5), game.AttemptsLeft)
			assert.Equal(t, 0, len(game.Guesses))
		})
	}
}

type gameAdderStub2 struct {
	err error
}

func (g gameAdderStub2) Add(game session.Game) error {
	return g.err
}
