package decimal

// Equal checks if 2 values are equal.
// It converts to `Decimal` before comparing which can be lossy for floating point numbers.
// The comparison is exact by value without any tolerance for floating point errors.
func Equal[A, B Number](a A, b B) bool {
	v1, v2 := New(a).full(), New(b).full()
	return v1.Negative == v2.Negative && v1.Integer == v2.Integer && v1.Fraction == v2.Fraction
}

// Compare compares 2 values.
// It converts to `Decimal` before comparing which can be lossy for floating point numbers.
func Compare[A, B Number](a A, b B) int {
	v1, v2 := New(a).full(), New(b).full()
	// Check for equality
	if v1.Negative == v2.Negative && v1.Integer == v2.Integer && v1.Fraction == v2.Fraction {
		return 0
	}
	// Check if a is negative and b is positive
	if v1.Negative && !v2.Negative {
		return -1
	}
	// Check if a is positive and b is negative
	if !v1.Negative && v2.Negative {
		return 1
	}
	// Get sign of both values (due to the previous two checks, they must have the same sign)
	sign := 1
	if v1.Negative && v2.Negative {
		sign = -1
	}
	// Check if the absolute value of a is smaller than the absolute value of b
	// The sign decides what that means, if positive the larger absolute value is greater, if negative the larger absolute value is smaller
	if v1.Integer < v2.Integer || (v1.Integer == v2.Integer && v1.Fraction < v2.Fraction) {
		return -1 * sign
	}
	return sign
}

// IsZero checks if a decimal value is zero.
func (d Decimal) IsZero() bool {
	return d.Integer == 0 && d.Fraction == 0
}
