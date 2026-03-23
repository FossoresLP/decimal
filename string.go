package decimal

import (
	"fmt"
)

const cutoff = ^uint64(0) / 10 // 1844674407370955161

// NewFromString parses a decimal value from a string.
// The string must contain just the number with no additional characters around it.
// It will parse at most 19 digits after the decimal point.
// The integer component must fit into an unsigned 64-bit integer.
func NewFromString(s string) (Decimal, error) {
	d := Zero()
	l := len(s)
	if l == 0 {
		return Zero(), fmt.Errorf("no number in string: %s", s)
	}
	pos := 0
	gotNum := false
	if s[0] == '-' {
		d.Negative = true
		pos = 1
	}

	for ; pos < l; pos++ {
		if s[pos] >= '0' && s[pos] <= '9' {
			if d.Integer >= cutoff {
				if d.Integer > cutoff || s[pos] > '5' {
					return Zero(), fmt.Errorf("value overflows unsigned 64-bit integer: %s", s)
				}
			}
			d.Integer = d.Integer*10 + uint64(s[pos]-'0')
			gotNum = true
		} else if s[pos] == '.' {
			pos++
			goto fracloop
		} else {
			return Zero(), fmt.Errorf("invalid character in integer: %s", s[pos:])
		}
	}
	if !gotNum {
		return Zero(), fmt.Errorf("no number in string: %s", s)
	}
	if d.Integer == 0 {
		d.Negative = false
	}
	return d, nil
fracloop:
	for ; pos < l; pos++ {
		if s[pos] >= '0' && s[pos] <= '9' {
			if d.Digits >= 19 {
				return Zero(), fmt.Errorf("more digits in fraction than can be represented: %d", pos)
			}
			d.Digits++
			d.Fraction = d.Fraction*10 + uint64(s[pos]-'0')
			gotNum = true
		} else {
			return Zero(), fmt.Errorf("invalid character in fraction: %s", s[pos:])
		}
	}
	if !gotNum {
		return Zero(), fmt.Errorf("no number in string: %s", s)
	}
	if d.Integer == 0 && d.Fraction == 0 {
		d.Negative = false
	}
	return d, nil
}

// NewFromStringFuzzy parses the first decimal value it can find from a string.
// Leading or trailing characters are ignored.
// The first digit, minus or dot marks the beginning of a number.
// It will parse at most 19 digits after the decimal point.
// The integer component must fit into an unsigned 64-bit integer.
func NewFromStringFuzzy(s string) (Decimal, error) {
	d := Zero()
	l := len(s)
	if l == 0 {
		return Zero(), fmt.Errorf("no digits found in string: %s", s)
	}
	pos := 0
numloop:
	for ; pos < l; pos++ {
		if s[pos] == '-' || (s[pos] >= '0' && s[pos] <= '9') || s[pos] == '.' {
			goto intloop
		}
	}
	return Zero(), fmt.Errorf("no digits found in string: %s", s)
intloop:
	if s[pos] == '-' {
		d.Negative = true
		pos++
	}
	gotNum := false
	for ; pos < l; pos++ {
		if s[pos] >= '0' && s[pos] <= '9' {
			if d.Integer >= cutoff {
				if d.Integer > cutoff || s[pos] > '5' {
					return Zero(), fmt.Errorf("value overflows unsigned 64-bit integer: %s", s)
				}
			}
			d.Integer = d.Integer*10 + uint64(s[pos]-'0')
			gotNum = true
		} else if s[pos] == '.' {
			pos++
			goto fracloop
		} else {
			if !gotNum {
				d = Zero()
				goto numloop
			}
			break
		}
	}
	if !gotNum {
		return Zero(), fmt.Errorf("no digits found in string: %s", s)
	}
	if d.Integer == 0 {
		d.Negative = false
	}
	return d, nil
fracloop:
	for ; pos < l; pos++ {
		if s[pos] >= '0' && s[pos] <= '9' {
			if d.Digits >= 19 {
				return Zero(), fmt.Errorf("more digits in fraction than can be represented: %d", pos)
			}
			d.Digits++
			d.Fraction = d.Fraction*10 + uint64(s[pos]-'0')
			gotNum = true
		} else {
			if !gotNum {
				d = Zero()
				goto numloop
			}
			break
		}
	}
	if !gotNum {
		return Zero(), fmt.Errorf("no digits found in string: %s", s)
	}
	if d.Integer == 0 && d.Fraction == 0 {
		d.Negative = false
	}
	return d, nil
}

// String converts a decimal value into a string representation.
func (d Decimal) String() string {
	arr := [48]byte{} // 1 digit sign, 20 digits integer, 1 dot, 20 digits fraction, aligned to 64-bit
	pos := 47
	if d.Digits > 0 {
		frac := d.Fraction
		end := 47 - int(d.Digits)
		for ; pos > end; pos-- {
			arr[pos] = byte(frac%10) + '0'
			frac /= 10
		}
		arr[pos] = '.'
		pos--
	}
	if d.Integer == 0 {
		arr[pos] = '0'
		pos--
	} else {
		for n := d.Integer; n > 0; n /= 10 {
			arr[pos] = byte(n%10) + '0'
			pos--
		}
	}
	if d.Negative {
		arr[pos] = '-'
	} else {
		pos++
	}
	return string(arr[pos:])
}
