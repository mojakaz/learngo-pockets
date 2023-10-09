package ecbank

import (
	"learngo-pockets/money"
	"testing"
)

func TestEnvelopeExchangeRate(t *testing.T) {
	tt := map[string]struct {
		source   string
		target   string
		envelope envelope
		expected money.ExchangeRate
		err      error
	}{
		"nominal": {
			source: "USD",
			target: "JPY",
			envelope: envelope{Rates: []currencyRate{
				{
					Currency: "USD",
					Rate:     1.02,
				},
				{
					Currency: "JPY",
					Rate:     145,
				},
			}},
			expected: money.ExchangeRate(142.156862745),
			err:      nil,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := tc.envelope.exchangeRate(tc.source, tc.target)
			if err != nil {
				t.Errorf("expected no error, got %v", err.Error())
			}
			if got != tc.expected {
				t.Errorf("expected %f, got %f", tc.expected, got)
			}
		})
	}

}
