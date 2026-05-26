package decimal

import (
	"math"
)

// Fixed is a fixed-point decimal number with 2 digits after the decimal point.
// It handles values up to [-21474836.48, 21474836.47].
type Fixed int32

// NewFixed converts any of the builtin numeric types in Go to a fixed-point decimal value.
// It handles conversions in the range [-21474836.48, 21474836.47], values outside that range overflow.
// Only two digits after the decimal point are considered, the rest is truncated silently.
func NewFixed[N Number](value N) Fixed {
	switch v := any(value).(type) {
	case uint:
		return Fixed(v) * 100
	case uint8:
		return Fixed(v) * 100
	case uint16:
		return Fixed(v) * 100
	case uint32:
		return Fixed(v) * 100
	case uint64:
		return Fixed(v) * 100
	case int:
		return Fixed(v) * 100
	case int8:
		return Fixed(v) * 100
	case int16:
		return Fixed(v) * 100
	case int32:
		return Fixed(v) * 100
	case int64:
		return Fixed(v) * 100
	case float32:
		if math.IsInf(float64(v), 0) || math.IsNaN(float64(v)) {
			return 0
		}
		f := Fixed(v * 100)
		return f
	case float64:
		if math.IsInf(v, 0) || math.IsNaN(v) {
			return 0
		}
		f := Fixed(v * 100)
		return f
	case Decimal:
		f := Fixed(v.Integer) * 100
		var sign Fixed = 1
		if v.Negative {
			sign = -1
		}
		switch v.Digits {
		case 0:
			return sign * f
		case 1:
			return sign * (f + Fixed(v.Fraction*10))
		case 2:
			return sign * (f + Fixed(v.Fraction))
		default:
			return sign * (f + Fixed(v.Fraction/pow10[v.Digits-2]))
		}
	case Fixed:
		return v
	default:
		panic("unsupported type: generics failed")
	}
}

// Decimal converts a fixed point value to a full decimal value losslessly.
func (f Fixed) Decimal() Decimal {
	var d Decimal
	val := int64(f)
	if f < 0 {
		d.Negative = true
		val = -val
	}
	d.Fraction = uint64(val % 100)
	d.Digits = 2
	d.Integer = uint64(val / 100)
	return d
}
