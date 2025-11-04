package decimal

var pow10 = [20]uint64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000, 100000000000, 1000000000000, 10000000000000, 100000000000000, 1000000000000000, 10000000000000000, 100000000000000000, 1000000000000000000, 10000000000000000000}

type Number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func New[N Number](value N) Decimal {
	var d Decimal
	switch v := any(value).(type) {
	case uint:
		d.Integer = uint64(v)
	case uint8:
		d.Integer = uint64(v)
	case uint16:
		d.Integer = uint64(v)
	case uint32:
		d.Integer = uint64(v)
	case uint64:
		d.Integer = v
	case int:
		if v < 0 {
			d.Negative = true
			v = -v
		}
		d.Integer = uint64(v)
	case int8:
		if v < 0 {
			d.Negative = true
			v = -v
		}
		d.Integer = uint64(v)
	case int16:
		if v < 0 {
			d.Negative = true
			v = -v
		}
		d.Integer = uint64(v)
	case int32:
		if v < 0 {
			d.Negative = true
			v = -v
		}
		d.Integer = uint64(v)
	case int64:
		if v < 0 {
			d.Negative = true
			v = -v
		}
		d.Integer = uint64(v)
	case float32:
		d.FromFloat64(float64(v))
	case float64:
		d.FromFloat64(v)
	}
	return d
}

func Zero() Decimal {
	return Decimal{}
}

type Decimal struct {
	Negative bool
	Digits   uint8 // Number of digits in Fraction, must be less or equal to 20, not all values with 20 digits can be represented due to limits of uint64
	Integer  uint64
	Fraction uint64
}

func (d Decimal) Equal(d2 Decimal) bool {
	return d.Negative == d2.Negative && d.Integer == d2.Integer && d.Fraction == d2.Fraction && d.Digits == d2.Digits
}

func (d Decimal) MultiplyUint64(u uint64) Decimal {
	d.Integer *= u
	d.Fraction *= u
	d.Integer += d.Fraction / pow10[d.Digits]
	d.Fraction %= pow10[d.Digits]
	return d
}

func (d Decimal) DivideUint64(u uint64) Decimal {
	d.Fraction += d.Integer % u * pow10[d.Digits]
	d.Integer /= u
	d.Fraction /= u
	return d
}

func (d Decimal) Add(d2 Decimal) Decimal {
	// Extend the number with less digits to match the other
	if d.Digits < d2.Digits {
		d = d.ToDigits(d2.Digits)
	} else if d.Digits > d2.Digits {
		d2 = d2.ToDigits(d.Digits)
	}
	if d.Negative == d2.Negative {
		d.Integer += d2.Integer
		d.Fraction += d2.Fraction // TODO: check for overflow
		d.Integer += d.Fraction / pow10[d.Digits]
		d.Fraction %= pow10[d.Digits]
	} else {
		if d.Integer > d2.Integer {
			d.Integer -= d2.Integer
			if d.Fraction < d2.Fraction {
				d.Fraction += pow10[d.Digits]
				d.Integer--
			}
			d.Fraction -= d2.Fraction
		} else {
			d.Integer = d2.Integer - d.Integer
			if d.Fraction > d2.Fraction {
				d.Fraction -= d2.Fraction
			} else {
				d.Fraction = d2.Fraction - d.Fraction
				d.Negative = !d.Negative
			}
		}
	}
	return d
}

func (d Decimal) ToDigits(digits uint8) Decimal {
	for ; d.Digits < digits; d.Digits++ {
		d.Fraction *= 10
	}
	for ; d.Digits > digits; d.Digits-- {
		d.Fraction /= 10
	}
	return d
}
