package money

import (
	"fmt"
)

// Amount defines a quantity of money in a given Currency.
type Amount struct {
	quantity Decimal
	currency Currency
}

func (a Amount) String() string {
	return a.quantity.String() + " " + a.currency.String()
}

// validate returns an error if and only if an Amount is unsafe to use.
func (a Amount) validate() error {
	const maxAmount = 12
	switch {
	case len(fmt.Sprintf("%d", a.quantity.subunits)) > maxAmount:
		return ErrTooLarge
	case a.quantity.precision > a.currency.precision:
		return ErrTooPrecise
	}

	return nil
}

const (
	// ErrTooPrecise is returned if the number is too precise for the currency.
	ErrTooPrecise = Error("quantity is too precise")
)

func NewAmount(quantity Decimal, currency Currency) (Amount, error) {
	switch {
	case quantity.precision > currency.precision:
		// In order to avoid converting 0.00001 cent, let's exit now.
		return Amount{}, ErrTooPrecise
	case quantity.precision < currency.precision:
		quantity.subunits *= pow10(currency.precision - quantity.precision)
		quantity.precision = currency.precision
	}

	quantity.precision = currency.precision

	return Amount{quantity: quantity, currency: currency}, nil
}
