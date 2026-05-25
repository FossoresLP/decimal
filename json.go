package decimal

import "unsafe"

// MarshalJSON encodes a decimal value as a JSON number.
func (d Decimal) MarshalJSON() ([]byte, error) {
	var arr [48]byte
	pos := d.text(&arr)
	b := make([]byte, 48-pos)
	copy(b, arr[pos:])
	return b, nil
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
