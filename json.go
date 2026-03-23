package decimal

import "unsafe"

// MarshalJSON encodes a decimal value as a JSON number.
func (d Decimal) MarshalJSON() ([]byte, error) {
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
	return out, nil
}

// UnmarshalJSON decodes a JSON number or string into a decimal value.
// Strings must be a plain number and may not contain any escaped or non-numeric characters.
// `null` is decoded as zero to ensure missing values do not stop decoding entirely.
func (d *Decimal) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && data[0] == 'n' && data[1] == 'u' && data[2] == 'l' && data[3] == 'l' {
		*d = Zero()
		return nil
	}
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	val, err := NewFromString(unsafe.String(unsafe.SliceData(data), len(data)))
	if err != nil {
		return err
	}
	*d = val
	return nil
}
