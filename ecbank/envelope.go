package ecbank

import (
	"encoding/xml"
	"fmt"
	"io"
	"learngo-pockets/money"
)

const (
	baseCurrencyCode      = "EUR"
	ErrUnexpectedFormat   = ecbankError("unexpected format")
	ErrChangeRateNotFound = ecbankError("change rate not found")
)

type envelope struct {
	Rates []currencyRate `xml:"Cube>Cube>Cube"`
}

type currencyRate struct {
	Currency string             `xml:"currency,attr"`
	Rate     money.ExchangeRate `xml:"rate,attr"`
}

// exchangeRates builds a map of all the supported exchange rates.
func (e envelope) exchangeRates() map[string]money.ExchangeRate {
	rates := make(map[string]money.ExchangeRate, len(e.Rates)+1)

	for _, c := range e.Rates {
		rates[c.Currency] = c.Rate
	}

	rates[baseCurrencyCode] = 1.

	return rates
}

// exchangeRate reads the change rate from the Envelope's contents.
func (e envelope) exchangeRate(source, target string) (money.ExchangeRate, error) {
	if source == target {
		return 1., nil
	}

	rates := e.exchangeRates()

	sourceFactor, sourceFound := rates[source]
	if !sourceFound {
		return 0, fmt.Errorf("failed to find the source currency %s", source)
	}

	targetFactor, targetFound := rates[target]
	if !targetFound {
		return 0, fmt.Errorf("failed to find target currency %s", target)
	}

	return targetFactor / sourceFactor, nil
}

// readRateFromResponse reads the response and returns money.ExchangeRate
func readRateFromResponse(source, target string, respBody io.Reader) (money.ExchangeRate, error) {
	// read the response
	decoder := xml.NewDecoder(respBody)

	var ecbMessage envelope
	err := decoder.Decode(&ecbMessage)
	if err != nil {
		return 0., fmt.Errorf("%w: %s", ErrUnexpectedFormat, err)
	}

	rate, err := ecbMessage.exchangeRate(source, target)
	if err != nil {
		return 0., fmt.Errorf("%w: %s", ErrChangeRateNotFound, err)
	}
	return rate, nil
}
