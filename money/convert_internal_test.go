package money

import (
	"reflect"
	"testing"
)

func TestMultiply_Success(t *testing.T) {
	tt := map[string]struct {
		d        Decimal
		r        ExchangeRate
		expected Decimal
		err      error
	}{
		"20.00 * 0.04 = 0.8": {
			d:        Decimal{2000, 2},
			r:        0.04,
			expected: Decimal{8, 1},
			err:      nil,
		},
		"1_200.10 * 1.41 = 1_692.141": {
			d:        Decimal{120010, 2},
			r:        1.41,
			expected: Decimal{1692141, 3},
			err:      nil,
		},
		"0.8 * 0.3 = 0.24": {
			d:        Decimal{8, 1},
			r:        0.3,
			expected: Decimal{24, 2},
			err:      nil,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := multiply(tc.d, tc.r)
			if err != nil {
				t.Errorf("expected no error, got %v", err.Error())
			}
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestApplyExchangeRate_Success(t *testing.T) {
	tt := map[string]struct {
		in       Amount
		target   Currency
		rate     ExchangeRate
		expected Amount
	}{
		"Amount(1.52) * rate(1) USD to EUR": {
			in: Amount{
				quantity: Decimal{152, 2},
				currency: Currency{
					code:      "USD",
					precision: 2,
				},
			},
			target: Currency{
				code:      "EUR",
				precision: 2,
			},
			rate: 1,
			expected: Amount{
				quantity: Decimal{152, 2},
				currency: Currency{
					code:      "EUR",
					precision: 2,
				},
			},
		},
		"Amount(2.50) * rate(4)": {
			in: Amount{
				quantity: Decimal{
					subunits:  250,
					precision: 2,
				},
				currency: Currency{
					code:      "JPN",
					precision: 2,
				},
			},
			target: Currency{
				code:      "USD",
				precision: 2,
			},
			rate: 4,
			expected: Amount{
				quantity: Decimal{
					subunits:  1000,
					precision: 2,
				},
				currency: Currency{
					code:      "USD",
					precision: 2,
				},
			},
		},
		"Amount(4) * rate(2.5)": {
			in: Amount{
				quantity: Decimal{
					subunits:  4,
					precision: 0,
				},
				currency: Currency{
					code:      "IRR",
					precision: 0,
				},
			},
			target: Currency{
				code:      "NIR",
				precision: 0,
			},
			rate: 2.5,
			expected: Amount{
				quantity: Decimal{
					subunits:  10,
					precision: 0,
				},
				currency: Currency{
					code:      "NIR",
					precision: 0,
				},
			},
		},
		"Amount(3.14) * rate(2.52678)": {
			in: Amount{
				quantity: Decimal{
					subunits:  314,
					precision: 2,
				},
				currency: Currency{
					code:      "USD",
					precision: 2,
				},
			},
			target: Currency{
				code:      "JPN",
				precision: 2,
			},
			rate: 2.52678,
			expected: Amount{
				quantity: Decimal{
					subunits:  793,
					precision: 2,
				},
				currency: Currency{
					code:      "JPN",
					precision: 2,
				},
			},
		},
		"Amount(1.1) * rate(10)": {
			in: Amount{
				quantity: Decimal{
					subunits:  11,
					precision: 1,
				},
				currency: Currency{
					code:      "CNY",
					precision: 1,
				},
			},
			target: Currency{
				code:      "VND",
				precision: 1,
			},
			rate: 10,
			expected: Amount{
				quantity: Decimal{
					subunits:  110,
					precision: 1,
				},
				currency: Currency{
					code:      "VND",
					precision: 1,
				},
			},
		},
		"Amount(1_000_000_000.01) * rate(2)": {
			in: Amount{
				quantity: Decimal{
					subunits:  1_000_000_00001,
					precision: 2,
				},
				currency: Currency{
					code:      "USD",
					precision: 2,
				},
			},
			target: Currency{
				code:      "JPN",
				precision: 2,
			},
			rate: 2,
			expected: Amount{
				quantity: Decimal{
					subunits:  2_000_000_00002,
					precision: 2,
				},
				currency: Currency{
					code:      "JPN",
					precision: 2,
				},
			},
		},
		"Amount(265_413.87) * rate(0.00051)": {
			in: Amount{
				quantity: Decimal{
					subunits:  265_41387,
					precision: 2,
				},
				currency: Currency{
					code:      "USD",
					precision: 2,
				},
			},
			target: Currency{
				code:      "JPN",
				precision: 2,
			},
			rate: 0.00051,
			expected: Amount{
				quantity: Decimal{
					subunits:  13536,
					precision: 2,
				},
				currency: Currency{
					code:      "JPN",
					precision: 2,
				},
			},
		},
		"Amount(265_413) * rate(1)": {
			in: Amount{
				quantity: Decimal{
					subunits:  265_413,
					precision: 0,
				},
				currency: Currency{
					code:      "IRR",
					precision: 0,
				},
			},
			target: Currency{
				code:      "IQD",
				precision: 3,
			},
			rate: 1,
			expected: Amount{
				quantity: Decimal{
					subunits:  265_413000,
					precision: 3,
				},
				currency: Currency{
					code:      "IQD",
					precision: 3,
				},
			},
		},
		"Amount(2) * rate(1.337)": {
			in: Amount{
				quantity: Decimal{
					subunits:  200,
					precision: 2,
				},
				currency: Currency{
					code:      "USD",
					precision: 2,
				},
			},
			target: Currency{
				code:      "XYZ",
				precision: 5,
			},
			rate: 1.337,
			expected: Amount{
				quantity: Decimal{
					subunits:  267400,
					precision: 5,
				},
				currency: Currency{
					code:      "XYZ",
					precision: 5,
				},
			},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := applyExchangeRate(tc.in, tc.target, tc.rate)
			if err != nil {
				t.Errorf("expected no error, got %v", err.Error())
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestApplyExchangeRate_Failure(t *testing.T) {
	tt := map[string]struct {
		a        Amount
		target   Currency
		rate     ExchangeRate
		expected Amount
		err      error
	}{
		"Amount(2) * rate(1.33 * 10^16)": {
			a: Amount{
				quantity: Decimal{
					subunits:  2,
					precision: 0,
				},
				currency: Currency{
					code:      "IRR",
					precision: 0,
				},
			},
			target: Currency{
				code:      "XYZ",
				precision: 5,
			},
			rate:     1.33e+16,
			expected: Amount{},
			err:      ErrTooLarge,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := applyExchangeRate(tc.a, tc.target, tc.rate)
			if err == nil {
				t.Errorf("expected %v, got nil", tc.err)
			}
		})
	}
}
