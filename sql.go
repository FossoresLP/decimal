package decimal

import (
	"database/sql/driver"
	"fmt"
	"unsafe"
)

// Scan converts SQL data into a decimal value.
// It handles textual representation as well as floating point and integer values.
// Decoding follows the standards set by `New` and `NewFromString` respectively.
func (d *Decimal) Scan(value any) (err error) {
	if value == nil {
		*d = Zero()
		return nil
	}
	switch v := value.(type) {
	case []byte:
		val, err := NewFromString(unsafe.String(unsafe.SliceData(v), len(v)))
		if err != nil {
			return err
		}
		*d = val
		return nil
	case string:
		val, err := NewFromString(v)
		if err != nil {
			return err
		}
		*d = val
		return nil
	case float64:
		*d = New(v)
		return nil
	case int64:
		if v < 0 {
			d.Negative = true
			v = -v
		} else {
			d.Negative = false
		}
		d.Integer = uint64(v)
		d.Fraction = 0
		d.Digits = 0
		return nil
	case uint64:
		d.Negative = false
		d.Integer = v
		d.Fraction = 0
		d.Digits = 0
		return nil
	default:
		return fmt.Errorf("invalid type for Decimal: %T", value)
	}
}

// Value encodes a decimal value for SQL.
// It uses a string representation that is widely compatible with most databases and data types.
func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}
