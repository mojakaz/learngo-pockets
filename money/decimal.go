package money

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	// ErrInvalidDecimal is returned if the decimal is malformed.
	ErrInvalidDecimal = Error("unable to convert the decimal")

	// ErrInvalidPrecision is returned if the length of the fraction part of the decimal is less than 0 or greater than 30.
	ErrInvalidPrecision = Error("precision must be less than 30")

	// ErrTooLarge is returned if the quantity is too large - this would cause floating point precision errors.
	ErrTooLarge = Error("quantity over 10^12 is too large")
)

// Decimal can represent a floating-point number with a fixed precision.
// example: 1.52 = 152 * 10^(-2) will be stored as {152, 2}
type Decimal struct {
	// subunits is the amount of subunits. Multiply it by the precision to get the real value
	subunits int64
	// Number of "subunits" in a unit, expressed as a power of 10.
	precision byte
}

func (d *Decimal) String() string {
	if d.precision == 0 {
		return fmt.Sprintf("%d", d.subunits)
	}
	frac := d.subunits % pow10(d.precision)
	integer := d.subunits / pow10(d.precision)
	decimalFormat := "%d.%0" + strconv.Itoa(int(d.precision)) + "d"
	return fmt.Sprintf(decimalFormat, integer, frac)
}
func (d *Decimal) simplify() {
	for d.subunits%10 == 0 && d.precision > 0 {
		d.subunits /= 10
		d.precision--
	}
}

// ParseDecimal converts a string into its Decimal representation.
// It assumes that there is up to one decimal separator, and that the separator is '.' (full stop character).
func ParseDecimal(value string) (Decimal, error) {
	// find the position of the . and split on it.
	intPart, fracPart, _ := strings.Cut(value, ".")

	// maxDecimal is the number of digits in a thousand billion.
	const maxDecimal = 12

	// check the length of the separated strings
	if len(intPart) > maxDecimal {
		return Decimal{}, ErrTooLarge
	}
	if len(fracPart) > 30 {
		return Decimal{}, ErrInvalidPrecision
	}

	precision := byte(len(fracPart))

	// convert the string without the . to an integer. This could fail
	subunits, err := strconv.ParseInt(intPart+fracPart, 10, 64)
	if err != nil {
		return Decimal{}, fmt.Errorf("%w: %s", ErrInvalidDecimal, err.Error())
	}

	// return the result
	return Decimal{subunits, precision}, nil
}

// pow10 is a quick implementation of how to raise 10 to a given power.
// It's optimised for small powers, and slow for unusually high powers.
func pow10(power byte) int64 {
	switch power {
	case 0:
		return 1
	case 1:
		return 10
	case 2:
		return 100
	case 3:
		return 1000
	default:
		return int64(math.Pow(10, float64(power)))
	}
}
