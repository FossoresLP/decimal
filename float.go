package decimal

// Float64 converts a decimal value to floating point.
// Precision loss is minimized but not all values can be represented exactly.
func (d Decimal) Float64() float64 {
	if d.Negative {
		return -float64(d.Integer) - float64(d.Fraction)/float64(pow10[d.Digits])
	}
	return float64(d.Integer) + float64(d.Fraction)/float64(pow10[d.Digits])

}
