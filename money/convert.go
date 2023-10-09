package money

import "fmt"

// Convert applies the change rate to convert an amount to a target currency.
func Convert(amount Amount, to Currency, rates exchangeRates) (Amount, error) {
	rate, err := rates.FetchExchangeRate(amount.currency, to)
	if err != nil {
		return Amount{}, fmt.Errorf("cannot get the exchange rate: %w", err)
	}
	// Convert to the target currency applying the fetched change rate.
	convertedValue, err := applyExchangeRate(amount, to, rate)
	if err != nil {
		return Amount{}, err
	}
	// Validate the converted amount is in the handled bounded range.
	err = convertedValue.validate()
	if err != nil {
		return Amount{}, err
	}

	return convertedValue, nil
}

// applyExchangeRate returns a new Amount representing the input Amount multiplied by the ExchangeRate.
// The precision of the returned value is that of the target Currency.
// This function does not guarantee that the output amount is supported.
func applyExchangeRate(in Amount, target Currency, rate ExchangeRate) (Amount, error) {
	converted, err := multiply(in.quantity, rate)
	if err != nil {
		return Amount{}, err
	}

	switch {
	case converted.precision > target.precision:
		converted.subunits = converted.subunits / pow10(converted.precision-target.precision)
	case converted.precision < target.precision:
		converted.subunits = converted.subunits * pow10(target.precision-converted.precision)
	}

	converted.precision = target.precision

	return Amount{converted, target}, nil
}

// multiply a Decimal with an ExchangeRate and returns the product
func multiply(d Decimal, r ExchangeRate) (Decimal, error) {
	// First, convert the ExchangeRate to a Decimal
	rate, err := ParseDecimal(fmt.Sprintf("%g", r))
	if err != nil {
		return Decimal{}, fmt.Errorf("%w: exchange rate is %f", ErrInvalidDecimal, r)
	}
	dec := Decimal{
		subunits:  d.subunits * rate.subunits,
		precision: d.precision + rate.precision,
	}
	// Let's clean the representation a bit. Remove trailing zeroes.
	dec.simplify()

	return dec, nil
}
