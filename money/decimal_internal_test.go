package money

import (
	"errors"
	"testing"
)

func TestParseDecimal(t *testing.T) {
	tt := map[string]struct {
		decimal  string
		expected Decimal
		err      error
	}{
		"2 decimal digits": {
			decimal:  "1.25",
			expected: Decimal{125, 2},
			err:      nil,
		},
		"no decimal digits": {
			decimal:  "3",
			expected: Decimal{3, 0},
			err:      nil,
		},
		"suffix 0 as decimal digits": {
			decimal:  "3.930",
			expected: Decimal{3930, 3},
			err:      nil,
		},
		"prefix 0 as decimal digits": {
			decimal:  "3.093",
			expected: Decimal{3093, 3},
			err:      nil,
		},
		"multiple of 10": {
			decimal:  "0000100",
			expected: Decimal{100, 0},
			err:      nil,
		},
		"invalid decimal part": {
			decimal:  "-103a",
			expected: Decimal{},
			err:      ErrInvalidDecimal,
		},
		"invalid precision": {
			decimal:  "204.11111111111111111111111111111111",
			expected: Decimal{},
			err:      ErrInvalidPrecision,
		},
		"not a number": {
			decimal:  "NaN",
			expected: Decimal{},
			err:      ErrInvalidDecimal,
		},
		"empty string": {
			decimal:  "",
			expected: Decimal{},
			err:      ErrInvalidDecimal,
		},
		"too large": {
			decimal:  "1234567890123",
			expected: Decimal{},
			err:      ErrTooLarge,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := ParseDecimal(tc.decimal)
			if !errors.Is(err, tc.err) {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
