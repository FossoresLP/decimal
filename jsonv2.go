//go:build goexperiment.jsonv2

package decimal

import (
	"encoding/json/jsontext"
	"fmt"
	"unsafe"
)

// MarshalJSONTo implements encoding/json/v2.MarshalerTo.
func (d Decimal) MarshalJSONTo(enc *jsontext.Encoder) error {
	val, err := d.MarshalJSON()
	if err != nil {
		return err
	}
	return enc.WriteValue(val)
}

// UnmarshalJSONFrom implements encoding/json/v2.UnmarshalerFrom.
func (d *Decimal) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}
	switch val.Kind() {
	case jsontext.KindNull:
		*d = Zero()
		return nil
	case jsontext.KindString:
		val = val[1 : len(val)-1] // strip quotes
		fallthrough
	case jsontext.KindNumber:
		parsed, err := NewFromString(unsafe.String(unsafe.SliceData(val), len(val)))
		if err != nil {
			return err
		}
		*d = parsed
		return nil
	default:
		return fmt.Errorf("decimal: unsupported JSON kind: %v", val.Kind())
	}
}
