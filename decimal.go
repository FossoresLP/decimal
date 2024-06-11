package decimal

import (
	"database/sql/driver"
	"fmt"
	"math"

	"github.com/fxamacker/cbor/v2"
)

var pow10 = [20]uint64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000, 100000000000, 1000000000000, 10000000000000, 100000000000000, 1000000000000000, 10000000000000000, 100000000000000000, 1000000000000000000, 10000000000000000000}

type Number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func New[N Number](value N) *Decimal {
	var d Decimal
	switch v := any(value).(type) {
	case uint:
		d.Integer = uint64(v)
	case uint8:
		d.Integer = uint64(v)
	case uint16:
		d.Integer = uint64(v)
	case uint32:
		d.Integer = uint64(v)
	case uint64:
		d.Integer = v
	case int:
		if v < 0 {
			d.Negative = true
			v = -v
		}
		d.Integer = uint64(v)
	case int8:
		if v < 0 {
			d.Negative = true
			v = -v
		}
		d.Integer = uint64(v)
	case int16:
		if v < 0 {
			d.Negative = true
			v = -v
		}
		d.Integer = uint64(v)
	case int32:
		if v < 0 {
			d.Negative = true
			v = -v
		}
		d.Integer = uint64(v)
	case int64:
		if v < 0 {
			d.Negative = true
			v = -v
		}
		d.Integer = uint64(v)
	case float32:
		d.FromFloat64(float64(v))
	case float64:
		d.FromFloat64(v)
	}
	return &d
}

func NewFromString(s string) (*Decimal, error) {
	var d Decimal
	err := d.FromString(s)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func Zero() *Decimal {
	return &Decimal{}
}

type Decimal struct {
	Negative bool
	Integer  uint64
	Fraction uint64
	Digits   uint8 // Number of digits in Fraction, must be less or equal to 20, not all values with 20 digits can be represented due to limits of uint64
}

func (d Decimal) Float64() float64 {
	if d.Negative {
		return -float64(d.Integer) - float64(d.Fraction)/float64(pow10[d.Digits])
	}
	return float64(d.Integer) + float64(d.Fraction)/float64(pow10[d.Digits])

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

func (d *Decimal) FromFloat64(f float64) {
	if f < 0 {
		d.Negative = true
		f = -f
	}
	d.Integer = uint64(f)
	fraction := math.Abs(f - float64(d.Integer))

	const maxDigits = 18 // Reasonable precision to handle float64 precision limits
	var digits uint8

	for i := 0; i < maxDigits; i++ {
		fraction *= 10
		digits++
		if fraction == math.Round(fraction) {
			break
		}
	}

	d.Fraction = uint64(fraction)
	d.Digits = digits
}

func (d *Decimal) Scan(value interface{}) error {
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

func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	return d.FromString(string(data))
}

func (d Decimal) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(d.Float64())
}

func (d *Decimal) UnmarshalCBOR(data []byte) error {
	var f float64
	if err := cbor.Unmarshal(data, &f); err != nil {
		return err
	}
	d.FromFloat64(f)
	return nil
}

func (d *Decimal) MultiplyUint64(u uint64) *Decimal {
	d.Integer *= u
	d.Fraction *= u
	d.Integer += d.Fraction / pow10[d.Digits]
	d.Fraction %= pow10[d.Digits]
	return d
}

func (d *Decimal) DivideUint64(u uint64) *Decimal {
	d.Fraction += d.Integer % u * pow10[d.Digits]
	d.Integer /= u
	d.Fraction /= u
	return d
}

func (d *Decimal) ToDigits(digits uint8) *Decimal {
	for ; d.Digits < digits; d.Digits++ {
		d.Fraction *= 10
	}
	for ; d.Digits > digits; d.Digits-- {
		d.Fraction /= 10
	}
	return d
}

func (d *Decimal) Clone() *Decimal {
	return &Decimal{
		Negative: d.Negative,
		Integer:  d.Integer,
		Fraction: d.Fraction,
		Digits:   d.Digits,
	}
}
