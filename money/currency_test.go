package money

import (
	"testing"
)

func TestParseCurrency_Success(t *testing.T) {
	tt := map[string]struct {
		code     string
		expected Currency
		err      error
	}{
		"integer IRR": {
			code: "IRR",
			expected: Currency{
				code:      "IRR",
				precision: 0,
			},
			err: nil,
		},
		"tenth CNY": {
			code: "CNY",
			expected: Currency{
				code:      "CNY",
				precision: 1,
			},
			err: nil,
		},
		"tenth VND": {
			code: "VND",
			expected: Currency{
				code:      "VND",
				precision: 1,
			},
			err: nil,
		},
		"thousandth BHD": {
			code: "BHD",
			expected: Currency{
				code:      "BHD",
				precision: 3,
			},
			err: nil,
		},
		"hundredth USD": {
			code: "USD",
			expected: Currency{
				code:      "USD",
				precision: 2,
			},
			err: nil,
		},
		"hundredth EUR": {
			code: "EUR",
			expected: Currency{
				code:      "EUR",
				precision: 2,
			},
			err: nil,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := ParseCurrency(tc.code)
			if err != nil {
				t.Errorf("expected no error, got %s", err.Error())
			}
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestParseCurrency_Failure(t *testing.T) {
	tt := map[string]struct {
		code     string
		expected Currency
		err      error
	}{
		"too long": {
			code:     "ZZZZ",
			expected: Currency{},
			err:      ErrInvalidCurrencyCode,
		},
		"too short": {
			code:     "AA",
			expected: Currency{},
			err:      ErrInvalidCurrencyCode,
		},
		"digits": {
			code:     "123",
			expected: Currency{},
			err:      ErrInvalidCurrencyCode,
		},
		"small case": {
			code:     "usd",
			expected: Currency{},
			err:      ErrInvalidCurrencyCode,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := ParseCurrency(tc.code)
			if err == nil {
				t.Errorf("expected an error %v, got %v", tc.err, err.Error())
			}
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
