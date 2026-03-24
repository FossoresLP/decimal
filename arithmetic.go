package decimal

import (
	"math/bits"
)

// Add adds two decimals together.
// Its overflow behavior matches that of integers in Go.
func Add[A, B Number](a A, b B) Decimal {
	v1, v2 := New(a), New(b)
	// Extend the number with less digits to match the other
	if v1.Digits < v2.Digits {
		v1 = v1.ToDigits(v2.Digits)
	} else if v1.Digits > v2.Digits {
		v2 = v2.ToDigits(v1.Digits)
	}
	if v1.Negative == v2.Negative {
		v1.Integer += v2.Integer
		sum, carry := bits.Add64(v1.Fraction, v2.Fraction, 0)
		quo, rem := bits.Div64(carry, sum, pow10[v1.Digits])
		v1.Integer += quo
		v1.Fraction = rem
	} else {
		// Determine which operand has larger magnitude
		if v1.Integer > v2.Integer || (v1.Integer == v2.Integer && v1.Fraction >= v2.Fraction) {
			// |d| >= |d2|, result keeps d's sign
			v1.Integer -= v2.Integer
			if v1.Fraction < v2.Fraction {
				v1.Fraction += pow10[v1.Digits]
				v1.Integer--
			}
			v1.Fraction -= v2.Fraction
		} else {
			// |d2| > |d|, result takes d2's sign
			v1.Integer = v2.Integer - v1.Integer
			if v2.Fraction < v1.Fraction {
				v2.Fraction += pow10[v1.Digits]
				v1.Integer--
			}
			v1.Fraction = v2.Fraction - v1.Fraction
			v1.Negative = v2.Negative
		}
		// Canonicalize zero
		if v1.Integer == 0 && v1.Fraction == 0 {
			v1.Negative = false
		}
	}
	return v1
}

// Subtract subtracts one decimal value from another.
// Its overflow behavior matches that of integers in Go.
func Subtract[A, B Number](a A, b B) Decimal {
	v1, v2 := New(a), New(b)
	v2.Negative = !v2.Negative
	return Add(v1, v2)
}

// Multiply multiplies two decimal values.
// Its overflow behavior matches that of integers in Go.
func Multiply[A, B Number](a A, b B) Decimal {
	v1, v2 := New(a), New(b)

	out := Decimal{
		Negative: v1.Negative != v2.Negative,   // Negative is XOR of both inputs
		Digits:   min(v1.Digits+v2.Digits, 19), // Output precision is the sum of input precisions, capped at 19
		Integer:  v1.Integer * v2.Integer,      // Term 1: Int × Int → pure integer (wraps on overflow)
	}

	// 128-bit fractional accumulator
	var fracHi, fracLo uint64

	// Term 2: Int1 × Frac2 → divide by 10^d2, rescale remainder to dOut
	if v1.Integer != 0 && v2.Fraction != 0 {
		hi, lo := bits.Mul64(v1.Integer, v2.Fraction)
		quo, rem := bits.Div64(hi, lo, pow10[v2.Digits])
		out.Integer += quo
		fracLo, fracHi = bits.Add64(fracLo, rem*pow10[out.Digits-v2.Digits], 0)
	}

	// Term 3: Frac1 × Int2 → divide by 10^d1, rescale remainder to dOut
	if v1.Fraction != 0 && v2.Integer != 0 {
		hi, lo := bits.Mul64(v1.Fraction, v2.Integer)
		quo, rem := bits.Div64(hi, lo, pow10[v1.Digits])
		out.Integer += quo
		var carry uint64
		fracLo, carry = bits.Add64(fracLo, rem*pow10[out.Digits-v1.Digits], 0)
		fracHi += carry
	}

	// Term 4: Frac1 × Frac2 → at scale 10^(d1+d2), truncate to dOut
	if v1.Fraction != 0 && v2.Fraction != 0 {
		hi, lo := bits.Mul64(v1.Fraction, v2.Fraction)
		var contrib uint64
		if excess := v1.Digits + v2.Digits - out.Digits; excess > 0 {
			contrib, _ = bits.Div64(hi, lo, pow10[excess])
		} else {
			contrib = lo
		}
		var carry uint64
		fracLo, carry = bits.Add64(fracLo, contrib, 0)
		fracHi += carry
	}

	// Carry from fraction into integer
	if fracHi > 0 || fracLo >= pow10[out.Digits] {
		quo, rem := bits.Div64(fracHi, fracLo, pow10[out.Digits])
		out.Integer += quo
		out.Fraction = rem
	} else {
		out.Fraction = fracLo
	}

	if out.Integer == 0 && out.Fraction == 0 {
		out.Negative = false
	}

	return out
}

// div128 divides a 128-bit numerator (numHi:numLo) by a 128-bit denominator (denHi:denLo).
// It returns the quotient as a uint64 (wrapping on overflow) and a 128-bit remainder.
func div128(numHi, numLo, denHi, denLo uint64) (quo, remHi, remLo uint64) {
	if denHi == 0 {
		if numHi < denLo {
			quo, remLo = bits.Div64(numHi, numLo, denLo)
		} else {
			r := numHi % denLo
			quo, remLo = bits.Div64(r, numLo, denLo)
		}
		return
	}

	// denHi > 0: quotient fits in 64 bits (num < 2^128, den >= 2^64)
	s := uint(bits.LeadingZeros64(denHi))
	if s == 0 {
		// den >= 2^127, quotient is 0 or 1
		if numHi > denHi || (numHi == denHi && numLo >= denLo) {
			quo = 1
			var borrow uint64
			remLo, borrow = bits.Sub64(numLo, denLo, 0)
			remHi, _ = bits.Sub64(numHi, denHi, borrow)
		} else {
			remHi, remLo = numHi, numLo
		}
		return
	}

	// Normalize: shift both sides so denHi's leading bit is set
	dd1 := (denHi << s) | (denLo >> (64 - s))
	dd0 := denLo << s
	nn2 := numHi >> (64 - s)
	nn1 := (numHi << s) | (numLo >> (64 - s))
	nn0 := numLo << s

	// Trial quotient from top 128 bits of shifted numerator / top 64 bits of shifted denominator
	quo, r := bits.Div64(nn2, nn1, dd1)

	// Knuth refinement: check q * dd0 against r:nn0
	pH, pL := bits.Mul64(quo, dd0)
	if pH > r || (pH == r && pL > nn0) {
		quo--
		r += dd1
		if r >= dd1 { // no overflow
			pH, pL = bits.Mul64(quo, dd0)
			if pH > r || (pH == r && pL > nn0) {
				quo--
			}
		}
	}

	// Remainder = num - quo * den (using original unshifted values)
	qLHi, qLLo := bits.Mul64(quo, denLo)
	_, qHLo := bits.Mul64(quo, denHi)
	mid, _ := bits.Add64(qHLo, qLHi, 0)
	var borrow uint64
	remLo, borrow = bits.Sub64(numLo, qLLo, 0)
	remHi, _ = bits.Sub64(numHi, mid, borrow)
	return
}

// Divide divides two decimal values.
// The result is computed with maximum precision (19 fractional digits).
// Its overflow and divide-by-zero behavior match that of integers in Go.
func Divide[A, B Number](a A, b B) Decimal {
	v1, v2 := New(a), New(b)

	if v2.Integer == 0 && v2.Fraction == 0 {
		panic("invalid operation: division by zero")
	}

	if v1.Integer == 0 && v1.Fraction == 0 {
		return Zero()
	}

	// Fast path when divisor is integer
	if v2.Fraction == 0 {
		v1.Fraction *= pow10[19-v1.Digits]
		v1.Digits = 19
		v1.Negative = v1.Negative != v2.Negative
		hi, lo := bits.Mul64(v1.Integer%v2.Integer, pow10[19])
		lo, carry := bits.Add64(lo, v1.Fraction, 0)
		hi += carry
		v1.Integer /= v2.Integer
		v1.Fraction, _ = bits.Div64(hi, lo, v2.Integer)
		v1.Integer += v1.Fraction / pow10[19]
		v1.Fraction %= pow10[19]
		if v1.Integer == 0 && v1.Fraction == 0 {
			v1.Negative = false
		}
		return v1.Truncate()
	}

	out := Decimal{
		Negative: v1.Negative != v2.Negative,
		Digits:   19,
	}

	// Scale both to 19-digit fixed-point (128-bit scaled integers)
	numHi, numLo := bits.Mul64(v1.Integer, pow10[19])
	f1 := v1.Fraction * pow10[19-v1.Digits]
	var c uint64
	numLo, c = bits.Add64(numLo, f1, 0)
	numHi += c

	denHi, denLo := bits.Mul64(v2.Integer, pow10[19])
	f2 := v2.Fraction * pow10[19-v2.Digits]
	denLo, c = bits.Add64(denLo, f2, 0)
	denHi += c

	// Integer quotient and remainder
	q, remHi, remLo := div128(numHi, numLo, denHi, denLo)
	out.Integer = q

	// Fractional quotient: rem * 10^19 / den
	if denHi == 0 {
		// Fast path: den fits in 64 bits, compute fraction in one step
		rHi, rLo := bits.Mul64(remLo, pow10[19])
		out.Fraction, _ = bits.Div64(rHi, rLo, denLo)
	} else {
		// General case: extract one fractional digit at a time
		var frac uint64
		for range 19 {
			// rem *= 10 (up to 192 bits: carry:remHi:remLo)
			pH, pL := bits.Mul64(remLo, 10)
			cH, cL := bits.Mul64(remHi, 10)
			remLo = pL
			var carry uint64
			remHi, carry = bits.Add64(cL, pH, 0)
			carry += cH

			// digit = carry:remHi:remLo / denHi:denLo (digit is at most 9)
			var digit uint64
			for carry > 0 || remHi > denHi || (remHi == denHi && remLo >= denLo) {
				var borrow uint64
				remLo, borrow = bits.Sub64(remLo, denLo, 0)
				remHi, borrow = bits.Sub64(remHi, denHi, borrow)
				carry -= borrow
				digit++
			}
			frac = frac*10 + digit
		}
		out.Fraction = frac
	}

	if out.Integer == 0 && out.Fraction == 0 {
		out.Negative = false
	}

	return out.Truncate()
}

// Negate returns the negation of the given value.
func Negate[N Number](n N) Decimal {
	v := New(n)
	if v.Integer != 0 || v.Fraction != 0 {
		v.Negative = !v.Negative
	}
	return v
}

// Absolute returns the absolute value of a given input.
func Absolute[N Number](n N) Decimal {
	v := New(n)
	v.Negative = false
	return v
}
