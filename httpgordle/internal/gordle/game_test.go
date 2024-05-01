package gordle

import (
	"testing"
)

func mustNewGame(corpus []string) *Game {
	game, _ := New(corpus)
	return game
}

func TestGame_Play(t *testing.T) {
	tt := map[string]struct {
		corpus   []string
		guess    string
		expected Feedback
	}{
		"correct": {
			corpus:   []string{"hello"},
			guess:    "hello",
			expected: Feedback{correctPosition, correctPosition, correctPosition, correctPosition, correctPosition},
		},
		"wrong position": {
			corpus:   []string{"hello"},
			guess:    "elohl",
			expected: Feedback{wrongPosition, wrongPosition, wrongPosition, wrongPosition, wrongPosition},
		},
		"absent character": {
			corpus:   []string{"hello"},
			guess:    "sprit",
			expected: Feedback{absentCharacter, absentCharacter, absentCharacter, absentCharacter, absentCharacter},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			game := mustNewGame(tc.corpus)
			got, err := game.Play(tc.guess)
			if err != nil {
				t.Errorf("expected no error, got an error: %s", err.Error())
			}
			if !got.Equal(tc.expected) {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestGame_ShowAnswer(t *testing.T) {
	tt := map[string]struct {
		corpus   []string
		expected string
	}{
		"correct": {
			corpus:   []string{"hello"},
			expected: "HELLO",
		},
		"wrong": {
			corpus:   []string{"wrong"},
			expected: "WRONG",
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			game := mustNewGame(tc.corpus)
			got := game.ShowAnswer()
			if got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestGame_ShowLength(t *testing.T) {
	tt := map[string]struct {
		corpus   []string
		expected int
	}{
		"correct": {
			corpus:   []string{"hello"},
			expected: 5,
		},
		"organization": {
			corpus:   []string{"organization"},
			expected: 12,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			game := mustNewGame(tc.corpus)
			got := game.ShowLength()
			if got != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, got)
			}
		})
	}
}
