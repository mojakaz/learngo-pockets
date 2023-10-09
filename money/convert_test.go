package money_test

import (
	"learngo-pockets/money"
	"reflect"
	"testing"
)

// stubRate is a very simple stub for the exchangeRates.
type stubRate struct {
	rate money.ExchangeRate
	err  error
}

// FetchExchangeRate implements the interface exchangeRates with the same signature but fields are unused for tests purposes.
func (m stubRate) FetchExchangeRate(_, _ money.Currency) (money.ExchangeRate, error) {
	return m.rate, m.err
}

func mustParseCurrency(t *testing.T, code string) money.Currency {
	t.Helper()
	currency, err := money.ParseCurrency(code)
	if err != nil {
		t.Fatalf("cannot parse currency %s code", code)
	}
	return currency
}

func mustParseAmount(t *testing.T, value, code string) money.Amount {
	t.Helper()
	decimal, err := money.ParseDecimal(value)
	if err != nil {
		t.Fatalf("invalid decimal value: %s", value)
	}
	currency, err := money.ParseCurrency(code)
	if err != nil {
		t.Fatalf("invalid currency code: %s", code)
	}
	amount, err := money.NewAmount(decimal, currency)
	if err != nil {
		t.Fatalf("cannot create Amount with value %v and code %s", decimal, code)
	}
	return amount
}

func TestConvert(t *testing.T) {
	tt := map[string]struct {
		// input    fields
		amount   money.Amount
		to       money.Currency
		rates    stubRate
		validate func(t *testing.T, got money.Amount, err error)
	}{
		"34.98 USD to EUR": {
			amount: mustParseAmount(t, "34.98", "USD"),
			to:     mustParseCurrency(t, "EUR"),
			rates: stubRate{
				rate: 2,
				err:  nil,
			},
			validate: func(t *testing.T, got money.Amount, err error) {
				if err != nil {
					t.Errorf("expected no error, got %s", err.Error())
				}
				expected := mustParseAmount(t, "69.96", "EUR")
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("expected %v, got %v", expected, got)
				}
			},
		},
		"29.3 JPY to USD": {
			amount: mustParseAmount(t, "29.3", "JPY"),
			to:     mustParseCurrency(t, "USD"),
			rates: stubRate{
				rate: 2,
				err:  nil,
			},
			validate: func(t *testing.T, got money.Amount, err error) {
				if err != nil {
					t.Errorf("expected no error, got %v", err.Error())
				}
				expected := mustParseAmount(t, "58.6", "USD")
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("expected %v, got %v", expected, got)
				}
			},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := money.Convert(tc.amount, tc.to, tc.rates)
			tc.validate(t, got, err)
		})
	}
}
