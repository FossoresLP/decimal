package decimal

import (
	"fmt"
)

func NewFromString(s string) (*Decimal, error) {
	var d Decimal
	err := d.FromString(s)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// NewFromStringFuzzy extracts and parses the first number from a string, ignoring any leading or trailing non-numeric characters.
func NewFromStringFuzzy(s string) (*Decimal, error) {
	var d Decimal
	err := d.FromStringFuzzy(s)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (d *Decimal) bytes() []byte {
	buf := [48]byte{} // need 48 bytes, 1 digit sign, 20 digits integer, 1 dot, 20 digits fraction, aligned to 64-bit
	pos := 47
	if d.Digits > 0 {
		frac := d.Fraction
		end := 47 - int(d.Digits)
		for ; pos > end; pos-- {
			buf[pos] = byte(frac%10) + '0'
			frac /= 10
		}
		buf[pos] = '.'
		pos--
	}
	if d.Integer == 0 {
		buf[pos] = '0'
		pos--
	} else {
		for n := d.Integer; n > 0; n /= 10 {
			buf[pos] = byte(n%10) + '0'
			pos--
		}
	}
	if d.Negative {
		buf[pos] = '-'
	} else {
		pos++
	}
	out := make([]byte, 48-pos)
	copy(out, buf[pos:])
	return out
}

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

func (d *Decimal) FromString(s string) error {
	d.Negative = false
	d.Integer = 0
	d.Fraction = 0
	d.Digits = 0
	l := len(s)
	if l == 0 {
		return nil
	}
	pos := 0
	if s[0] == '-' {
		d.Negative = true
		pos = 1
	}
	for ; pos < l; pos++ {
		if s[pos] >= '0' && s[pos] <= '9' {
			d.Integer = d.Integer*10 + uint64(s[pos]-'0')
		} else if s[pos] == '.' {
			pos++
			goto fracloop
		} else {
			return fmt.Errorf("invalid character in integer: %s", s[pos:])
		}
	}
	return nil
fracloop:
	for ; pos < l; pos++ {
		if s[pos] >= '0' && s[pos] <= '9' {
			d.Digits++
			d.Fraction = d.Fraction*10 + uint64(s[pos]-'0')
		} else {
			return fmt.Errorf("invalid character in fraction: %s", s[pos:])
		}
	}
	return nil
}

// FromStringFuzzy extracts and parses the first number from a string, ignoring any leading or trailing non-numeric characters.
func (d *Decimal) FromStringFuzzy(s string) error {
	d.Negative = false
	d.Integer = 0
	d.Fraction = 0
	d.Digits = 0
	l := len(s)
	if l == 0 {
		return nil
	}
	pos := 0
	for ; pos < l; pos++ {
		if s[pos] == '-' || s[pos] >= '0' && s[pos] <= '9' || s[pos] == '.' {
			goto intloop
		}
	}
	return fmt.Errorf("no digits found in string: %s", s)
intloop:
	if s[0] == '-' {
		d.Negative = true
		pos = 1
	}
	for ; pos < l; pos++ {
		if s[pos] >= '0' && s[pos] <= '9' {
			d.Integer = d.Integer*10 + uint64(s[pos]-'0')
		} else if s[pos] == '.' {
			pos++
			goto fracloop
		} else {
			return nil
		}
	}
	return nil
fracloop:
	for ; pos < l; pos++ {
		if s[pos] >= '0' && s[pos] <= '9' {
			d.Digits++
			d.Fraction = d.Fraction*10 + uint64(s[pos]-'0')
		} else {
			return nil
		}
	}
	return nil
}
