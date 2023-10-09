package money

import (
	"fmt"
	"regexp"
)

// Currency defines the code of a money and its decimal precision.
type Currency struct {
	code      string
	precision byte
}

// String implements the Stringer interface.
func (c Currency) String() string {
	return c.code
}

func (c Currency) Code() string {
	return c.code
}

// ErrInvalidCurrencyCode is returned when the currency to parse is not a standard 3-letter code.
const ErrInvalidCurrencyCode = Error("invalid currency code")

// ParseCurrency returns the currency associated to a name and may return ErrInvalidCurrencyCode.
func ParseCurrency(code string) (Currency, error) {
	if len(code) != 3 {
		return Currency{}, ErrInvalidCurrencyCode
	}
	matched, err := regexp.MatchString("[A-Z]{3}", code)
	if err != nil {
		return Currency{}, fmt.Errorf("%w: %v", ErrInvalidCurrencyCode, err.Error())
	}
	if !matched {
		return Currency{}, fmt.Errorf("%w: currency code must be made of 3 letters between A and Z", ErrInvalidCurrencyCode)
	}

	switch code {
	case "IRR":
		return Currency{code: code, precision: 0}, nil
	case "CNY", "VND":
		return Currency{code: code, precision: 1}, nil
	case "BHD", "IQD", "KWD", "LYD", "OMR", "TND":
		return Currency{code: code, precision: 3}, nil
	default:
		return Currency{code: code, precision: 2}, nil
	}
}
