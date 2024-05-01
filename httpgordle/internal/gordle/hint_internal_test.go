package gordle

import (
	"fmt"
	"testing"
)

func TestFeedback_String(t *testing.T) {
	tt := map[string]struct {
		feedback Feedback
		expected string
	}{
		"nominal": {
			feedback: Feedback{absentCharacter, wrongPosition, correctPosition},
			expected: "⬜️🟡💚",
		},
		"not exist": {
			feedback: Feedback{100},
			expected: "💔",
		},
		"mix": {
			feedback: Feedback{3, absentCharacter, wrongPosition, correctPosition, 8},
			expected: "💔⬜️🟡💚💔",
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := fmt.Sprint(tc.feedback)
			if got != tc.expected {
				t.Errorf("different feedback: got %s, want %s", got, tc.expected)
			}
		})
	}
}

func TestFeedback_GameWon(t *testing.T) {
	tt := map[string]struct {
		feedback Feedback
		expected bool
	}{
		"nominal gameWon=false": {
			feedback: Feedback{absentCharacter, wrongPosition, correctPosition},
			expected: false,
		},
		"nominal gameWon=true": {
			feedback: Feedback{correctPosition, correctPosition, correctPosition},
			expected: true,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := tc.feedback.GameWon()
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
