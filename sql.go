package decimal

import (
	"database/sql/driver"
	"fmt"
)

func (d *Decimal) Scan(value any) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return d.FromString(string(v))
	case string:
		return d.FromString(v)
	case float64:
		d.FromFloat64(v)
		return nil
	case int64:
		if v < 0 {
			d.Negative = true
			v = -v
		}
		d.Integer = uint64(v)
		d.Fraction = 0
		d.Digits = 0
		return nil
	case uint64:
		d.Integer = v
		d.Fraction = 0
		d.Digits = 0
		return nil
	default:
		return fmt.Errorf("invalid type for Decimal: %T", value)
	}
}

func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}
