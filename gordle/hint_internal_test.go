package gordle

import (
	"fmt"
	"testing"
)

func TestFeedback_String(t *testing.T) {
	tt := map[string]struct {
		feedback feedback
		expected string
	}{
		"nominal": {
			feedback: []hint{absentCharacter, wrongPosition, correctPosition},
			expected: "⬜️🟡💚",
		},
		"not exist": {
			feedback: []hint{100},
			expected: "💔",
		},
		"mix": {
			feedback: []hint{3, absentCharacter, wrongPosition, correctPosition, 8},
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
