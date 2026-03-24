package decimal

import (
	"math"
	"math/bits"
)

var pow10 = [20]uint64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000, 100000000000, 1000000000000, 10000000000000, 100000000000000, 1000000000000000, 10000000000000000, 100000000000000000, 1000000000000000000, 10000000000000000000}

type Number interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64 | Decimal
}

// New converts any of the builtin numeric types in Go to a decimal value.
// It represents integer values exactly.
// 32-bit floating point numbers are converted with 10 digits behind the decimal point.
// 64-bit floating point numbers are converted with 18 digits behind the decimal point.
// Floating point numbers exceeding the unsigned 64-bit integer range inherit Go's float-to-uint64 overflow logic.
// Positive and negative infinity and NaN cannot be represented and are instead converted to zero.
func New[N Number](value N) Decimal {
	switch v := any(value).(type) {
	case uint:
		return Decimal{
			Integer: uint64(v),
		}
	case uint8:
		return Decimal{
			Integer: uint64(v),
		}
	case uint16:
		return Decimal{
			Integer: uint64(v),
		}
	case uint32:
		return Decimal{
			Integer: uint64(v),
		}
	case uint64:
		return Decimal{
			Integer: uint64(v),
		}
	case int:
		var d Decimal
		i := int64(v)
		if i < 0 {
			d.Negative = true
			i = -i
		}
		d.Integer = uint64(i)
		return d
	case int8:
		var d Decimal
		i := int64(v)
		if i < 0 {
			d.Negative = true
			i = -i
		}
		d.Integer = uint64(i)
		return d
	case int16:
		var d Decimal
		i := int64(v)
		if i < 0 {
			d.Negative = true
			i = -i
		}
		d.Integer = uint64(i)
		return d
	case int32:
		var d Decimal
		i := int64(v)
		if i < 0 {
			d.Negative = true
			i = -i
		}
		d.Integer = uint64(i)
		return d
	case int64:
		var d Decimal
		if v < 0 {
			d.Negative = true
			v = -v
		}
		d.Integer = uint64(v)
		return d
	case float32:
		var d Decimal
		if math.IsInf(float64(v), 0) || math.IsNaN(float64(v)) {
			return Zero()
		}
		d.Negative = v < 0
		if d.Negative {
			v = -v
		}
		d.Integer = uint64(v)
		d.Fraction = uint64((v - float32(d.Integer)) * float32(pow10[10]))
		d.Digits = 10
		d = d.Truncate()
		return d
	case float64:
		var d Decimal
		if math.IsInf(v, 0) || math.IsNaN(v) {
			return Zero()
		}
		d.Negative = v < 0
		if d.Negative {
			v = -v
		}
		d.Integer = uint64(v)
		d.Fraction = uint64((v - float64(d.Integer)) * float64(pow10[18]))
		d.Digits = 18
		d = d.Truncate()
		return d
	case Decimal:
		return v
	default:
		panic("unsupported type: generics failed")
	}
}

// Zero returns a zero value decimal
func Zero() Decimal {
	return Decimal{}
}

// Decimal represents high precision decimal numbers.
// It can store numbers with up to 19 digits behind the decimal point.
// The integer component is constrained to the range of a 64-bit unsigned integer in both positive and negative values.
// When manually constructed or modified, `Digits` must represent the exact number of digits in `Fraction`.
// Values with more than 19 fractional digits as well as negative zero are unsupported and will trigger undefined behavior.
type Decimal struct {
	Negative bool
	Digits   uint8 // Number of digits in Fraction, must be less or equal to 19
	Integer  uint64
	Fraction uint64
}

func (d Decimal) full() Decimal {
	d.Fraction *= pow10[19-d.Digits]
	return d
}

// ToDigits converts a decimal value to the specified number of digits after the decimal point.
// The number of digits is limited to 19.
// Digits beyond the defined number are truncated, no rounding is performed.
func (d Decimal) ToDigits(digits uint8) Decimal {
	if digits > 19 {
		digits = 19
	}
	for ; d.Digits < digits; d.Digits++ {
		d.Fraction *= 10
	}
	for ; d.Digits > digits; d.Digits-- {
		d.Fraction /= 10
	}
	if d.Integer == 0 && d.Fraction == 0 {
		d.Negative = false
	}
	return d
}

// Truncate removes trailing zeros from the decimal value.
func (d Decimal) Truncate() Decimal {
	if d.Fraction == 0 {
		if d.Integer == 0 {
			d.Negative = false
		}
		d.Digits = 0
		return d
	}

	if quot, rem := bits.Div64(0, d.Fraction, 10000000000000000); rem == 0 { // Check for 16 zeros
		d.Fraction = quot
		d.Digits -= 16
	}
	if quot, rem := bits.Div64(0, d.Fraction, 100000000); rem == 0 { // Check for 8 zeros
		d.Fraction = quot
		d.Digits -= 8
	}
	if quot, rem := bits.Div64(0, d.Fraction, 10000); rem == 0 { // Check for 4 zeros
		d.Fraction = quot
		d.Digits -= 4
	}
	if quot, rem := bits.Div64(0, d.Fraction, 100); rem == 0 { // Check for 2 zeros
		d.Fraction = quot
		d.Digits -= 2
	}
	if quot, rem := bits.Div64(0, d.Fraction, 10); rem == 0 { // Check for 1 zero
		d.Fraction = quot
		d.Digits -= 1
	}

	return d
}
