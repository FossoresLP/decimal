package decimal

func (d Decimal) Float64() float64 {
	if d.Negative {
		return -float64(d.Integer) - float64(d.Fraction)/float64(pow10[d.Digits])
	}
	return float64(d.Integer) + float64(d.Fraction)/float64(pow10[d.Digits])

}

// FromFloat64 converts a float64 to a decimal with up to 18 digits after the decimal point, avoiding trailing zeros.
// Digits beyond the 18th digit after the decimal point are truncated and values above MaxInt64 and below MinInt64 can't be converted. NaN and Inf are not supported.
func (d *Decimal) FromFloat64(f float64) {
	d.FromFloat64Fixed(f, 18)
	if d.Fraction == 0 {
		d.Digits = 0
		return
	}
	for d.Fraction%10 == 0 {
		d.Fraction /= 10
		d.Digits--
	}
}

// FromFloat64Fixed converts a float64 to a decimal with a fixed precision, truncating any additional digits and potentially leading to trailing zeros.
func (d *Decimal) FromFloat64Fixed(f float64, digits uint8) {
	if digits > 18 {
		digits = 18
	}
	d.Negative = f < 0
	if d.Negative {
		f = -f
	}
	d.Integer = uint64(f)
	d.Fraction = uint64((f - float64(d.Integer)) * float64(pow10[digits]))
	d.Digits = digits
}
