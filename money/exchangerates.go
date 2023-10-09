package money

type exchangeRates interface {
	FetchExchangeRate(source, target Currency) (ExchangeRate, error)
}

// ExchangeRate represents a rate to convert from a currency to another.
type ExchangeRate float32
