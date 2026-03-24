package decimal

import (
	"fmt"
	"io"
	"unsafe"
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

// Format implements fmt.Formatter.
// It supports "%f" and "%v" but not the other formats specified for floating point values such as "%e".
func (d Decimal) Format(state fmt.State, verb rune) {
	pad := func(char byte, count int) {
		if count <= 0 {
			return
		}
		buf := [64]byte{}
		for i := range buf {
			buf[i] = char
		}
		for count > 0 {
			chunk := min(count, len(buf))
			state.Write(buf[:chunk])
			count -= chunk
		}
	}

	switch verb {
	case 'f':
		trail := 0
		if precision, ok := state.Precision(); ok {
			d = d.Round(uint8(max(0, min(precision, 19))))
			if precision > int(d.Digits) {
				trail = min(precision, 1024) - int(d.Digits)
			}
		} else {
			d = d.Round(6)
		}
		width, fixedWidth := state.Width()
		width = min(width, 1024)
		sign := ""
		if d.Negative {
			sign = "-"
		} else if state.Flag('+') {
			sign = "+"
		} else if state.Flag(' ') {
			sign = " "
		}
		d.Negative = false
		str := d.String()
		if d.Digits == 0 && state.Flag('#') {
			str = str + "."
		}
		length := len(sign) + len(str) + trail
		if !fixedWidth {
			io.WriteString(state, sign)
			io.WriteString(state, str)
			pad('0', trail)
			return
		}
		if state.Flag('-') {
			io.WriteString(state, sign)
			io.WriteString(state, str)
			pad('0', trail)
			pad(' ', width-length)
			return
		}
		if state.Flag('0') {
			io.WriteString(state, sign)
			pad('0', width-length)
			io.WriteString(state, str)
			pad('0', trail)
			return
		}
		pad(' ', width-length)
		io.WriteString(state, sign)
		io.WriteString(state, str)
		pad('0', trail)
	case 'v':
		// %v handles only padding via state.Width() and state.Flag('-') and debug formatting via state.Flag('#').
		// Other flags are intentially ignored since even the standard library has an inconsistent approach for those.
		width, fixedWidth := state.Width()
		left := state.Flag('-')

		str := d.String()
		if state.Flag('#') {
			str = fmt.Sprintf("decimal.Decimal{Negative: %t, Integer: %d, Fraction: %d, Digits: %d}", d.Negative, d.Integer, d.Fraction, d.Digits)
		}

		if fixedWidth && !left {
			pad(' ', width-len(str))
		}
		io.WriteString(state, str)
		if fixedWidth && left {
			pad(' ', width-len(str))
		}
	default:
		fmt.Fprintf(state, "%%!%c(decimal.Decimal=%s)", verb, d.String())
	}
}

// MarshalText implements encoding.TextMarshaler.
func (d Decimal) MarshalText() ([]byte, error) {
	return d.MarshalJSON()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *Decimal) UnmarshalText(data []byte) error {
	val, err := NewFromString(unsafe.String(unsafe.SliceData(data), len(data)))
	if err != nil {
		return err
	}
	*d = val
	return nil
}
