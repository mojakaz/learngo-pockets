package money

import (
	"errors"
	"testing"
)

func TestNewAmount(t *testing.T) {
	tt := map[string]struct {
		quantity Decimal
		currency Currency
		expected Amount
		err      error
	}{
		"nominal": {
			quantity: Decimal{
				subunits:  102,
				precision: 2,
			},
			currency: Currency{
				code:      "USD",
				precision: 2,
			},
			expected: Amount{
				quantity: Decimal{
					subunits:  102,
					precision: 2,
				},
				currency: Currency{
					code:      "USD",
					precision: 2,
				},
			},
			err: nil,
		},
		"nominal2": {
			quantity: Decimal{
				subunits:  125,
				precision: 1,
			},
			currency: Currency{
				code:      "BHD",
				precision: 3,
			},
			expected: Amount{
				quantity: Decimal{
					subunits:  12500,
					precision: 3,
				},
				currency: Currency{
					code:      "BHD",
					precision: 3,
				},
			},
			err: nil,
		},
		"too precise": {
			quantity: Decimal{
				subunits:  234,
				precision: 2,
			},
			currency: Currency{
				code:      "CNY",
				precision: 1,
			},
			expected: Amount{},
			err:      ErrTooPrecise,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := NewAmount(tc.quantity, tc.currency)
			if err != nil && tc.err == nil {
				t.Errorf("expected no error, got %s", err)
			} else if !errors.Is(err, tc.err) {
				t.Errorf("expected %s, got %s", tc.err, err)
			}

			if !EqualAmount(got, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func EqualAmount(got, expected Amount) bool {
	if got.quantity != expected.quantity {
		return false
	}
	if got.currency != expected.currency {
		return false
	}
	return true
}
