package decimal

import (
	"database/sql/driver"
	"fmt"
	"math"
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

// Scan converts SQL data into a fixed-point value.
// It handles textual representation as well as floating point and integer values.
// Textual representation follows the semantics of NewFixedFromString.
// Floating point errors are accounted for while genuine sub-hundredth precision is rejected.
// Integer values are taken exactly. Values outside the fixed-point range are rejected.
func (f *Fixed) Scan(value any) (err error) {
	if value == nil {
		*f = 0
		return nil
	}
	switch v := value.(type) {
	case []byte:
		val, err := NewFixedFromString(unsafe.String(unsafe.SliceData(v), len(v)))
		if err != nil {
			return err
		}
		*f = val
		return nil
	case string:
		val, err := NewFixedFromString(v)
		if err != nil {
			return err
		}
		*f = val
		return nil
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return fmt.Errorf("value out of range: %v", v)
		}
		scaled := v * 100
		cents := math.Round(scaled)

		if cents < math.MinInt32 || cents > math.MaxInt32 {
			return fmt.Errorf("value out of range: %v", v)
		}
		// Float noise in scaled is ~|scaled|*2^-53 (<1e-6 across the whole int32 range)
		if math.Abs(scaled-cents) > 1e-6 {
			return fmt.Errorf("too many decimal digits: %v", v)
		}

		*f = Fixed(cents)
		return nil
	case int64:
		if v < math.MinInt32/100 || v > math.MaxInt32/100 {
			return fmt.Errorf("value out of range: %d", v)
		}
		*f = Fixed(v * 100)
		return nil
	case uint64:
		if v > math.MaxInt32/100 {
			return fmt.Errorf("value out of range: %d", v)
		}
		*f = Fixed(v * 100)
		return nil
	default:
		return fmt.Errorf("invalid type for Fixed: %T", value)
	}
}

// Value encodes a fixed-point value for SQL.
// It uses a string representation that is widely compatible with most databases and data types.
func (f Fixed) Value() (driver.Value, error) {
	return f.String(), nil
}
