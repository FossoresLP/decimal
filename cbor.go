package decimal

import "github.com/fxamacker/cbor/v2"

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
