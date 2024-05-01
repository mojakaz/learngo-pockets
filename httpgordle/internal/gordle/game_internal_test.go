package gordle

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func TestSplitToUppercaseCharacters(t *testing.T) {
	tt := map[string]struct {
		input    string
		expected []rune
	}{
		"nominal": {
			input:    "input",
			expected: []rune("INPUT"),
		},
		"nominal2": {
			input:    "hello",
			expected: []rune("HELLO"),
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := splitToUppercaseCharacters(tc.input)
			if !slices.Equal(got, tc.expected) {
				t.Errorf("expected %c, got %c", tc.expected, got)
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
			expected: ErrInvalidGuess,
		},
		"too short": {
			word:     []rune("LOG"),
			expected: ErrInvalidGuess,
		},
		"blank": {
			word:     []rune(""),
			expected: ErrInvalidGuess,
		},
		"nil": {
			word:     nil,
			expected: ErrInvalidGuess,
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

func TestComputeFeedback_Success(t *testing.T) {
	tt := map[string]struct {
		guess    []rune
		solution []rune
		expected Feedback
	}{
		"nominal": {
			guess:    []rune("HEL"),
			solution: []rune("HEL"),
			expected: Feedback{correctPosition, correctPosition, correctPosition},
		},
		"one absent character": {
			guess:    []rune("AIM"),
			solution: []rune("AIN"),
			expected: Feedback{correctPosition, correctPosition, absentCharacter},
		},
		"two absent characters": {
			guess:    []rune("WHOA"),
			solution: []rune("IMOA"),
			expected: Feedback{absentCharacter, absentCharacter, correctPosition, correctPosition},
		},
		"one wrong position": {
			guess:    []rune("FUJI"),
			solution: []rune("COIN"),
			expected: Feedback{absentCharacter, absentCharacter, absentCharacter, wrongPosition},
		},
		"two wrong position": {
			guess:    []rune("OIMOI"),
			solution: []rune("IMTQS"),
			expected: Feedback{absentCharacter, wrongPosition, wrongPosition, absentCharacter, absentCharacter},
		},
		"three wrong position and one correct": {
			guess:    []rune("HELLO"),
			solution: []rune("LLEKO"),
			expected: Feedback{absentCharacter, wrongPosition, wrongPosition, wrongPosition, correctPosition},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := computeFeedback(tc.guess, tc.solution)
			if err != nil {
				t.Errorf("expected no error, got an error: %s", err.Error())
			}
			if !got.Equal(tc.expected) {
				t.Errorf("different feedback: got %s, want %s", got, tc.expected)
			}
		})
	}
}

func TestComputeFeedback_Failure(t *testing.T) {
	tt := map[string]struct {
		guess    []rune
		solution []rune
		expected Feedback
	}{
		"nominal": {
			guess:    []rune("HELLO"),
			solution: []rune("HEL"),
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := computeFeedback(tc.guess, tc.solution)
			assert.ErrorIs(t, err, ErrInvalidGuess, "expected an error, got no error")
		})
	}
}
