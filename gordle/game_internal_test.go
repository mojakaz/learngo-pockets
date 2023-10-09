package gordle

import (
	"errors"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestGameAsk(t *testing.T) {
	tt := map[string]struct {
		input string
		want  []rune
	}{
		"5 characters in english": {
			input: "HELLO",
			want:  []rune("HELLO"),
		},
		"5 characters in arabic": {
			input: "مرحبا",
			want:  []rune("مرحبا"),
		},
		"5 characters in japanese": {
			input: "こんにちは",
			want:  []rune("こんにちは"),
		},
		"3 characters in japanese": {
			input: "おはよ",
			want:  []rune("おはよ"),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			g, _ := New(strings.NewReader(tc.input), []string{string(tc.want)}, 5)

			got := g.ask()
			if !slices.Equal(got, tc.want) {
				t.Errorf("got = %v, want %v", string(got), string(tc.want))
			}
		})
	}
}

func TestGameValidateGuess(t *testing.T) {
	tt := map[string]struct {
		word     []rune
		expected error
	}{
		"nominal": {
			word:     []rune("GUESS"),
			expected: nil,
		},
		"too long": {
			word:     []rune("POCKET"),
			expected: errInvalidWordLength,
		},
		"too short": {
			word:     []rune("LOG"),
			expected: errInvalidWordLength,
		},
		"blank": {
			word:     []rune(""),
			expected: errInvalidWordLength,
		},
		"nil": {
			word:     nil,
			expected: errInvalidWordLength,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := validateGuess(tc.word, []rune("GUESS"))
			if !errors.Is(err, tc.expected) {
				t.Errorf("%c, expected %q, got %q", tc.word, tc.expected, err)
			}
		})
	}
}

func TestComputeFeedback(t *testing.T) {
	tt := map[string]struct {
		guess    []rune
		solution []rune
		expected feedback
	}{
		"nominal": {
			guess:    []rune("HEL"),
			solution: []rune("HEL"),
			expected: feedback{correctPosition, correctPosition, correctPosition},
		},
		"one absent character": {
			guess:    []rune("AIM"),
			solution: []rune("AIN"),
			expected: feedback{correctPosition, correctPosition, absentCharacter},
		},
		"two absent characters": {
			guess:    []rune("WHOA"),
			solution: []rune("IMOA"),
			expected: feedback{absentCharacter, absentCharacter, correctPosition, correctPosition},
		},
		"one wrong position": {
			guess:    []rune("FUJI"),
			solution: []rune("COIN"),
			expected: feedback{absentCharacter, absentCharacter, absentCharacter, wrongPosition},
		},
		"two wrong position": {
			guess:    []rune("OIMOI"),
			solution: []rune("IMTQS"),
			expected: feedback{absentCharacter, wrongPosition, wrongPosition, absentCharacter, absentCharacter},
		},
		"three wrong position and one correct": {
			guess:    []rune("HELLO"),
			solution: []rune("LLEKO"),
			expected: feedback{absentCharacter, wrongPosition, wrongPosition, wrongPosition, correctPosition},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := computeFeedback(tc.guess, tc.solution)
			if !got.Equal(tc.expected) {
				t.Errorf("different feedback: got %s, want %s", got, tc.expected)
			}
		})
	}
}
