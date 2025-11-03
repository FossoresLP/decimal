package decimal

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"

	"github.com/x448/float16"
)

const (
	CBOR_MAJOR           = 0b111_00000
	CBOR_ADDITIONAL      = 0b000_11111
	CBOR_INTPOS          = 0b000_00000
	CBOR_INTNEG          = 0b001_00000
	CBOR_BYTESTRING      = 0b010_00000
	CBOR_ARRAY_LEN2      = 0b100_00010
	CBOR_TAG             = 0b110_00000
	CBOR_TAG_BIGNUMPOS   = 0b110_00010
	CBOR_TAG_BIGNUMNEG   = 0b110_00011
	CBOR_TAG_DECIMALFRAC = 0b110_00100
	CBOR_TYPE7           = 0b111_00000
	CBOR_FLOAT16         = 0b111_11001
	CBOR_FLOAT32         = 0b111_11010
	CBOR_FLOAT64         = 0b111_11011
)

func cborPutInt(n uint64, buf []byte) int {
	if n <= 23 {
		buf[0] |= byte(n)
		return 1
	}
	if n <= math.MaxUint8 {
		buf[0] |= byte(24)
		buf[1] = byte(n)
		return 2
	}
	if n <= math.MaxUint16 {
		buf[0] |= byte(25)
		binary.BigEndian.PutUint16(buf[1:], uint16(n))
		return 3
	}
	if n <= math.MaxUint32 {
		buf[0] |= byte(26)
		binary.BigEndian.PutUint32(buf[1:], uint32(n))
		return 5
	}
	buf[0] |= byte(27)
	binary.BigEndian.PutUint64(buf[1:], n)
	return 9
}

// MarshalCBOR implements the cbor.Marshaler interface.
// It encodes the decimal number according to RFC 8949 Section 3.4.4 for Decimal Fractions.
// The format is a CBOR tag 4 containing a two-element array: [exponent, mantissa].
func (d Decimal) MarshalCBOR() ([]byte, error) {
	out := make([]byte, 24)

	if d.Digits == 0 {
		if d.Negative {
			out[0] = CBOR_INTNEG
			n := cborPutInt(d.Integer-1, out)
			return out[:n], nil
		}
		out[0] = CBOR_INTPOS
		n := cborPutInt(d.Integer, out)
		return out[:n], nil
	}

	out[0] = CBOR_TAG_DECIMALFRAC
	out[1] = CBOR_ARRAY_LEN2

	if d.Digits >= uint8(len(pow10)) {
		return nil, fmt.Errorf("decimal: Digits (%d) out of range for encoding", d.Digits)
	}

	out[2] = CBOR_INTNEG | (d.Digits - 1)

	hi, lo := bits.Mul64(d.Integer, pow10[d.Digits])
	lo, carry := bits.Add64(lo, d.Fraction, 0)
	hi += carry

	if hi == 0 {
		if d.Negative {
			out[3] = CBOR_INTNEG
			n := cborPutInt(lo-1, out[3:])
			return out[:3+n], nil
		}
		out[3] = CBOR_INTPOS
		n := cborPutInt(lo, out[3:])
		return out[:3+n], nil
	}

	if d.Negative {
		out[3] = CBOR_TAG_BIGNUMNEG
		lo, borrow := bits.Sub64(lo, 1, 0)
		hi -= borrow
		bytes := 16 - bits.LeadingZeros64(hi)/8
		out[4] = CBOR_BYTESTRING | byte(bytes)
		binary.BigEndian.PutUint64(out[8:], hi)
		binary.BigEndian.PutUint64(out[16:], lo)
		copy(out[5:], out[23-bytes:])
		return out[:4+bytes], nil
	}
	out[3] = CBOR_TAG_BIGNUMPOS
	bytes := 16 - bits.LeadingZeros64(hi)/8
	out[4] = CBOR_BYTESTRING | byte(bytes)
	binary.BigEndian.PutUint64(out[8:], hi)
	binary.BigEndian.PutUint64(out[16:], lo)
	copy(out[5:], out[23-bytes:])
	return out[:4+bytes], nil
}

func cborParseInt(buf []byte) (uint64, int, bool, error) {
	if len(buf) < 1 {
		return 0, 0, false, fmt.Errorf("cbor: not enough data for integer at %d bytes", len(buf))
	}
	neg := false
	bytes := 0
	switch buf[0] & CBOR_MAJOR {
	case CBOR_INTPOS:
		neg = false
	case CBOR_INTNEG:
		neg = true
	default:
		return 0, 0, false, fmt.Errorf("cbor: unexpected major type %d when parsing integer", buf[0]>>5)
	}
	var val uint64
	additional := buf[0] & CBOR_ADDITIONAL
	switch additional {
	case 24:
		if len(buf) < 2 {
			return 0, 0, false, fmt.Errorf("cbor: not enough data for integer with additional information 24 at %d bytes", len(buf))
		}
		val = uint64(buf[1])
		bytes = 2
	case 25:
		if len(buf) < 3 {
			return 0, 0, false, fmt.Errorf("cbor: not enough data for integer with additional information 25 at %d bytes", len(buf))
		}
		val = uint64(binary.BigEndian.Uint16(buf[1:]))
		bytes = 3
	case 26:
		if len(buf) < 5 {
			return 0, 0, false, fmt.Errorf("cbor: not enough data for integer with additional information 26 at %d bytes", len(buf))
		}
		val = uint64(binary.BigEndian.Uint32(buf[1:]))
		bytes = 5
	case 27:
		if len(buf) < 9 {
			return 0, 0, false, fmt.Errorf("cbor: not enough data for integer with additional information 27 at %d bytes", len(buf))
		}
		val = binary.BigEndian.Uint64(buf[1:])
		bytes = 9
	default:
		if additional >= 24 {
			return 0, 0, false, fmt.Errorf("cbor: unexpected additional information %d when parsing integer", additional)
		}
		val = uint64(additional)
		bytes = 1
	}
	if neg {
		if val == math.MaxUint64 {
			return 0, 0, false, fmt.Errorf("cbor: negative integer overflows uint64")
		}
		val++
	}
	return val, bytes, neg, nil
}

// UnmarshalCBOR implements the cbor.Unmarshaler interface.
// It supports decoding RFC 8949 Decimal Fractions, standard floats, and integers.
// This implementation uses a fast path for mantissas that fit in a 64-bit integer
// and a slower path using math/big for larger numbers.
func (d *Decimal) UnmarshalCBOR(data []byte) error {
	if len(data) < 1 {
		return fmt.Errorf("cbor: input too short at %d bytes", len(data))
	}
	switch data[0] & CBOR_MAJOR {
	case CBOR_INTPOS, CBOR_INTNEG:
		d.Digits = 0
		d.Fraction = 0
		val, bytes, neg, err := cborParseInt(data)
		if err != nil {
			return err
		}
		if bytes != len(data) {
			return fmt.Errorf("cbor: integer did not consume entire data")
		}
		d.Integer = val
		d.Negative = neg
		return nil
	case CBOR_TAG:
		switch data[0] {
		case CBOR_TAG_DECIMALFRAC:
			if len(data) < 4 {
				return fmt.Errorf("cbor: not enough data for decimal fraction")
			}
			if data[1] != CBOR_ARRAY_LEN2 {
				return fmt.Errorf("cbor: decimal fraction tag not followed by two-element array")
			}
			expVal, expBytes, expNeg, err := cborParseInt(data[2:])
			if err != nil {
				return err
			}
			if expVal > 19 {
				return fmt.Errorf("cbor: decimal fraction power exceeds uint64")
			}
			if len(data) < 2+expBytes+1 {
				return fmt.Errorf("cbor: not enough data for mantissa in decimal fraction")
			}
			switch data[2+expBytes] & CBOR_MAJOR {
			case CBOR_INTPOS, CBOR_INTNEG:
				mantVal, mantBytes, mantNeg, err := cborParseInt(data[2+expBytes:])
				if err != nil {
					return err
				}
				if 2+expBytes+mantBytes < len(data) {
					return fmt.Errorf("cbor: decimal fraction did not consume entire data")
				}
				d.Negative = mantNeg
				if !expNeg {
					hi, lo := bits.Mul64(mantVal, pow10[expVal])
					if hi != 0 {
						return fmt.Errorf("cbor: decimal fraction value overflows uint64")
					}
					d.Digits = 0
					d.Fraction = 0
					d.Integer = lo
					return nil
				}
				d.Integer, d.Fraction = bits.Div64(0, mantVal, pow10[expVal])
				d.Digits = uint8(expVal)
				return nil
			case CBOR_TAG:
				mantNeg := false
				switch data[2+expBytes] & CBOR_ADDITIONAL {
				case CBOR_TAG_BIGNUMNEG:
					mantNeg = true
					fallthrough
				case CBOR_TAG_BIGNUMPOS:
					if len(data) < 2+expBytes+2 {
						return fmt.Errorf("cbor: not enough data for mantissa of decimal fraction")
					}
					if data[2+expBytes+1]&CBOR_MAJOR != CBOR_BYTESTRING {
						return fmt.Errorf("cbor: bignum tag is not followed by bytestring")
					}
					mantBytes := int(data[2+expBytes+1] & CBOR_ADDITIONAL)
					if mantBytes > 16 {
						return fmt.Errorf("cbor: bignum for mantissa is too large for decimal value")
					}
					if len(data) < 2+expBytes+2+mantBytes {
						return fmt.Errorf("cbor: not enough data for mantissa in decimal fraction, expected %d, got %d", mantBytes, len(data)-2-expBytes-2)
					}
					buf := [16]byte{}
					copy(buf[16-mantBytes:], data[2+expBytes+2:])
					hi := binary.BigEndian.Uint64(buf[:8])
					lo := binary.BigEndian.Uint64(buf[8:])
					if !expNeg {
						if hi != 0 {
							return fmt.Errorf("cbor: decimal fraction overflows uint64")
						}
						overflow, val := bits.Mul64(lo, expVal)
						if overflow != 0 {
							return fmt.Errorf("cbor: decimal fraction overflows uint64")
						}
						d.Digits = 0
						d.Fraction = 0
						d.Integer = val
						d.Negative = mantNeg
						return nil
					}
					if expVal <= hi {
						return fmt.Errorf("cbor: decimal fraction overflows uint64")
					}
					d.Integer, d.Fraction = bits.Div64(hi, lo, expVal)
					d.Digits = uint8(expVal)
					d.Negative = mantNeg
					return nil
				default:
					return fmt.Errorf("cbor: unexpected tag %d as mantissa in decimal fraction", data[2+expBytes]&CBOR_ADDITIONAL)
				}
			default:
				return fmt.Errorf("cbor: unexpected major type %d as mantissa in decimal fraction", data[2+expBytes]>>5)
			}
		case CBOR_TAG_BIGNUMPOS, CBOR_TAG_BIGNUMNEG:
			return fmt.Errorf("cbor: cannot parse bignum as decimal as it would overflow uint64")
		default:
			return fmt.Errorf("cbor: unexpected tag %d, cannot parse as decimal", data[0]&CBOR_ADDITIONAL)
		}
	case CBOR_TYPE7:
		switch data[0] {
		case CBOR_FLOAT16:
			if len(data) != 3 {
				return fmt.Errorf("cbor: unexpected length %d for float16", len(data))
			}
			d.FromFloat64(float64(float16.Frombits(binary.BigEndian.Uint16(data[1:])).Float32()))
			return nil
		case CBOR_FLOAT32:
			if len(data) != 5 {
				return fmt.Errorf("cbor: unexpected length %d for float32", len(data))
			}
			d.FromFloat64(float64(math.Float32frombits(binary.BigEndian.Uint32(data[1:]))))
			return nil
		case CBOR_FLOAT64:
			if len(data) != 9 {
				return fmt.Errorf("cbor: unexpected length %d for float64", len(data))
			}
			d.FromFloat64(math.Float64frombits(binary.BigEndian.Uint64(data[1:])))
			return nil
		default:
			return fmt.Errorf("cbor: unexpected simple value %d, cannot parse as decimal", data[0]&CBOR_ADDITIONAL)
		}
	default:
		return fmt.Errorf("cbor: unexpected major type %d, cannot parse as decimal", data[0]>>5)
	}
}
