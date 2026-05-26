package decimal

import (
	"fmt"
	"io"
	"math"
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

// internal helper for text conversion.
// 1 digit sign, 20 digits integer, 1 dot, 20 digits fraction, aligned to 64-bit
func (d Decimal) text(arr *[48]byte) int {
	pos := 48
	if d.Digits > 0 {
		for ; d.Fraction > 0; d.Fraction /= 10 {
			pos--
			arr[pos] = byte(d.Fraction%10) + '0'
		}
		frac := 48 - int(d.Digits)
		for pos > frac {
			pos--
			arr[pos] = '0'
		}
		pos--
		arr[pos] = '.'
	}
	pos--
	arr[pos] = byte(d.Integer%10) + '0'
	d.Integer /= 10
	for ; d.Integer > 0; d.Integer /= 10 {
		pos--
		arr[pos] = byte(d.Integer%10) + '0'
	}
	if d.Negative {
		pos--
		arr[pos] = '-'
	}
	return pos
}

// String converts a decimal value into a string representation.
func (d Decimal) String() string {
	var arr [48]byte
	pos := d.text(&arr)
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
	var arr [48]byte
	pos := d.text(&arr)
	b := make([]byte, 48-pos)
	copy(b, arr[pos:])
	return b, nil
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

// NewFixedFromString parses a fixed-point value from a string.
// The string must contain just the number with no additional characters around it.
// Fractional digits that cannot be represented are rejected.
// Strings with more than 16 characters plus optional sign are rejected outright.
// The value must fit in the range [-21474836.48, 21474836.47].
func NewFixedFromString(s string) (Fixed, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("no number in string: %s", s)
	}
	var sign Fixed = 1
	if s[0] == '-' {
		sign = -1
		s = s[1:]
	}
	if len(s) == 0 {
		return 0, fmt.Errorf("no number in string: %s", s)
	}
	if len(s) > 16 {
		return 0, fmt.Errorf("value too long to be processed as fixed decimal: %s", s)
	}
	var val uint64 = 0
	var hasFrac bool = false
	for i := range s {
		if c := s[i] - '0'; c <= 9 {
			val = val*10 + uint64(c)
		} else if s[i] == '.' {
			if i == 0 && len(s) == 1 {
				return 0, fmt.Errorf("no number in string: %s", s)
			}
			hasFrac = true
			s = s[i+1:]
			break
		} else {
			return 0, fmt.Errorf("invalid character in integer: %c", s[i])
		}
	}
	val *= 100
	if !hasFrac {
		if val > math.MaxInt32 {
			return 0, fmt.Errorf("value overflows fixed decimal: %d", val)
		}
		return sign * Fixed(val), nil
	}
	var f1, f2 uint64
	if len(s) > 0 {
		f1 = uint64(s[0] - '0')
		if f1 > 9 {
			return 0, fmt.Errorf("invalid character in fraction: %c", s[0])
		}
	}
	if len(s) > 1 {
		f2 = uint64(s[1] - '0')
		if f2 > 9 {
			return 0, fmt.Errorf("invalid character in fraction: %c", s[1])
		}
	}
	if len(s) > 2 {
		for i := range s[2:] {
			if s[i+2] != '0' {
				return 0, fmt.Errorf("invalid character in fraction overflow: %c", s[i+2])
			}
		}
	}
	val += f1*10 + f2

	if val > math.MaxInt32 {
		if sign < 0 && val == math.MaxInt32+1 {
			return Fixed(math.MinInt32), nil
		}
		return 0, fmt.Errorf("value overflows fixed decimal: %d", val)
	}
	return sign * Fixed(val), nil
}

// internal helper for text conversion.
// 1 digit sign, 8 digits integer, 1 dot, 2 digits fraction, aligned to 64-bit
func (f Fixed) text(arr *[16]byte) int {
	val := int64(f)
	neg := val < 0
	if neg {
		val = -val
	}
	arr[15] = byte(val%10) + '0'
	val /= 10
	arr[14] = byte(val%10) + '0'
	val /= 10
	arr[13] = '.'
	arr[12] = byte(val%10) + '0'
	val /= 10
	pos := 12
	for ; val > 0; val /= 10 {
		pos--
		arr[pos] = byte(val%10) + '0'
	}
	if neg {
		pos--
		arr[pos] = '-'
	}
	return pos
}

// String converts a fixed-point decimal value into a string representation.
func (f Fixed) String() string {
	var arr [16]byte
	pos := f.text(&arr)
	return string(arr[pos:])
}

// MarshalText implements encoding.TextMarshaler.
func (f Fixed) MarshalText() ([]byte, error) {
	var arr [16]byte
	pos := f.text(&arr)
	b := make([]byte, 16-pos)
	copy(b, arr[pos:])
	return b, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (f *Fixed) UnmarshalText(data []byte) error {
	val, err := NewFixedFromString(unsafe.String(unsafe.SliceData(data), len(data)))
	if err != nil {
		return err
	}
	*f = val
	return nil
}
